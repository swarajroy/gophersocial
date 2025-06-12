package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrUnauthenticated   = errors.New("unauthenticated")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleID    int64    `json:"-"`
	Role      Role     `json:"-"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) IsValid(input string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(input))
}

func (p *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	p.text = &text
	p.hash = hash
	return nil
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {

	query := `INSERT INTO users (username, email, password, role_id) 
	VALUES ($1, $2, $3, (SELECT id from roles where name = $4)) RETURNING id, username, created_at`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	err := tx.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password.hash,
		user.Role.Name,
	).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique key constraint (users_email_key)`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique key constraint (users_username_key)`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (us *UserStore) GetById(ctx context.Context, userID int64) (*User, error) {
	query := `SELECT u.id, u.username, u.email, u.password, u.created_at, u.is_active,
	r.id, r.name, r.level, r.description
	FROM 
	users u JOIN roles r ON (u.role_id = r.id)
	WHERE u.id = $1 and u.is_active = true`

	ctx, cancel := context.WithTimeout(ctx, QUERY_READ_TIME_OUR_DURATION)
	defer cancel()

	var user = &User{}

	err := us.db.QueryRowContext(ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.IsActive,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (us *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {

	query := `SELECT id, username, email, password, created_at 
	FROM 
	users 
	WHERE email = $1 and is_active = true`

	ctx, cancel := context.WithTimeout(ctx, QUERY_READ_TIME_OUR_DURATION)
	defer cancel()

	var user = &User{}

	err := us.db.QueryRowContext(ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (u *UserStore) createUserInvitation(ctx context.Context, tx *sql.Tx, user *User, token string, invitationExp time.Duration) error {

	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, user.ID, time.Now().Add(invitationExp))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	//with transaction wrapper
	return WithTx(u.db, ctx, func(tx *sql.Tx) error {
		// create user
		if err := u.Create(ctx, tx, user); err != nil {
			return err
		}
		// create invitation
		if err := u.createUserInvitation(ctx, tx, user, token, invitationExp); err != nil {
			return err
		}

		return nil
	})
}

func (u *UserStore) Activate(ctx context.Context, token string) error {
	//with transaction wrapper
	return WithTx(u.db, ctx, func(tx *sql.Tx) error {
		// 1. find the user from user invitations to which this token belongs
		user, err := u.getUserFromUserInvitation(ctx, tx, token)
		if err != nil {
			return err
		}
		// 2. update the user is_active to true from false
		user.IsActive = true
		if err := u.update(ctx, tx, user); err != nil {
			return err
		}
		// 3. clean up the invitation for that user from user_invitations
		if err := u.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}
		return nil
	})
}

func (u *UserStore) getUserFromUserInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.is_active 
	FROM users u 
	JOIN user_invitations ui ON u.id = ui.user_id 
	WHERE
	ui.token = $1 AND ui.expiry > $2
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	user := &User{}

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (u *UserStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET username = $1, email = $2, is_active = $3 WHERE id = $4`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return nil
		}
	}
	return nil
}

func (u *UserStore) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE from user_invitations where user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return nil
		}
	}
	return nil
}

func (us *UserStore) Delete(ctx context.Context, userID int64) error {
	WithTx(us.db, ctx, func(tx *sql.Tx) error {
		if err := us.delete(ctx, tx, userID); err != nil {
			return err
		}

		if err := us.deleteUserInvitation(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
	return nil
}

func (us *UserStore) delete(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE from users where user_id = $1`

	_, err := tx.ExecContext(ctx, query, userID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return ErrNotFound
		default:
			return err
		}
	}
	return nil
}

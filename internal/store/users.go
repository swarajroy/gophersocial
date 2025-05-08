package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateUsername = errors.New("duplicate username")
)

type User struct {
	ID        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  password `json:"-"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
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

	query := `INSERT INTO users (username, email, password) 
	VALUES ($1, $2, $3) RETURNING id, username, created_at`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	err := tx.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password.hash,
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
	query := `SELECT id, username, email, password, created_at 
	FROM 
	users 
	WHERE id = $1`

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
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, err
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

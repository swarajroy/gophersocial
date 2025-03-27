package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (us *UserStore) Create(ctx context.Context, user *User) error {

	query := `INSERT INTO users (username, email, password) 
	VALUES ($1, $2, $3) RETURNING id, username, created_at`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	err := us.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		return err
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

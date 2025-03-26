package store

import (
	"context"
	"database/sql"
)

type Followers struct {
	db *sql.DB
}

func NewFollowers(db *sql.DB) *Followers {
	return &Followers{
		db: db,
	}
}

func (f *Followers) Follow(ctx context.Context, followerId, userId int64) error {

	query := `
		INSERT into followers (user_id, follower_id)
		VALUES ($1, $2)
	`
	_, err := f.db.ExecContext(ctx, query, followerId, userId)
	return err
}

func (f *Followers) Unfollow(ctx context.Context, followerId, userId int64) error {
	query := `
		DELETE FROM followers 
		WHERE 
		user_id = $1 AND follower_id = $2
	`
	_, err := f.db.ExecContext(ctx, query, followerId, userId)
	return err
}

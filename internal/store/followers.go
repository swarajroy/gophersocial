package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
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
	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, followerId, userId)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return nil
}

func (f *Followers) Unfollow(ctx context.Context, followerId, userId int64) error {
	query := `
		DELETE FROM followers 
		WHERE 
		user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_READ_TIME_OUR_DURATION)
	defer cancel()

	_, err := f.db.ExecContext(ctx, query, followerId, userId)
	return err
}

package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const (
	QUERY_READ_TIME_OUR_DURATION  = time.Second * 3
	QUERY_WRITE_TIME_OUR_DURATION = time.Second * 5
)

var (
	ErrNotFound = errors.New("record not found")
	ErrConflict = errors.New("resource already exists")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedQuery) ([]PostWithMetadata, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetById(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetPostById(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followerId, userId int64) error
		Unfollow(ctx context.Context, followerId, userId int64) error
	}
	Roles interface {
		GetRoleByName(ctx context.Context, name string) (*Role, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     NewPostStore(db),
		Users:     NewUserStore(db),
		Comments:  NewCommentStore(db),
		Followers: NewFollowers(db),
		Roles:     NewRoleStore(db),
	}
}

func WithTx(db *sql.DB, ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

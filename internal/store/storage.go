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
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetPostById(context.Context, int64) ([]Comment, error)
	}
}

func NewStorage(db *sql.DB) Storage {
	return Storage{
		Posts:    NewPostStore(db),
		Users:    NewUserStore(db),
		Comments: NewCommentStore(db),
	}
}

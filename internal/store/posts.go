package store

import (
	"context"
	"database/sql"
)

type PostsStore struct {
	db *sql.DB
}

func NewPostsStore(db *sql.DB) *PostsStore {
	return &PostsStore{db: db}
}

func (ps *PostsStore) Create(ctx context.Context) error {
	return nil
}

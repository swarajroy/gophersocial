package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

func NewUsersStore(db *sql.DB) *UsersStore {
	return &UsersStore{db: db}
}

func (us *UsersStore) Create(ctx context.Context) error {
	return nil
}

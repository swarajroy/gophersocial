package store

import (
	"context"
	"database/sql"
	"time"
)

func NewMockDbStorage() Storage {
	return Storage{
		Users: &mockUserStore{},
	}
}

type mockUserStore struct {
}

func (us *mockUserStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	return nil
}

func (mu *mockUserStore) GetById(ctx context.Context, userID int64) (*User, error) {
	return nil, nil
}

func (us *mockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	return nil, nil
}

func (us *mockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return nil
}

func (us *mockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (us *mockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}

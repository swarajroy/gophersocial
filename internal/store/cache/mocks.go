package cache

import (
	"context"

	"github.com/swarajroy/gophersocial/internal/store"
)

func NewMockCacheStorage() Storage {
	return Storage{
		Users: &mockCacheUserStore{},
	}
}

type mockCacheUserStore struct {
}

func (m *mockCacheUserStore) Get(context.Context, int64) (*store.User, error) {
	return nil, nil
}

func (m *mockCacheUserStore) Set(context.Context, *store.User) error {
	return nil
}

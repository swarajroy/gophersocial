package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/swarajroy/gophersocial/internal/store"
)

const (
	UserCacheKeyEx = time.Duration(time.Minute)
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
)

type UserStore struct {
	rdb *redis.Client
}

func NewUserStore(rdb *redis.Client) *UserStore {
	return &UserStore{
		rdb: rdb,
	}
}

func (us *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)
	data, err := us.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}

	return &user, nil
}

func (us *UserStore) Set(ctx context.Context, user *store.User) error {

	if user.ID == 0 {
		return ErrInvalidUserID
	}

	cacheKey := fmt.Sprintf("user-%d", user.ID)

	data, err := json.Marshal(user)

	if err != nil {
		return err
	}

	return us.rdb.SetEX(ctx, cacheKey, data, UserCacheKeyEx).Err()
}

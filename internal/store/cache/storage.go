package cache

import (
	"context"

	"github.com/Aiyanu/gophersocial/internal/store"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}
}

func NewRedisStorage(rbd *redis.Client) Storage {
	return Storage{
		Users: &UserStore{rbd},
	}
}

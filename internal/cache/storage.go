package cache

import (
	"context"

	"github.com/deriannavy/microgo/internal/store"
	"github.com/go-redis/redis/v8"
)

type Storage struct {
	Account interface {
		Get(context.Context, int64) (*store.Account, error)
		Set(context.Context, *store.Account) error
		Delete(context.Context, int64)
	}
}

func NewCacheStorage(rbd *redis.Client) Storage {
	return Storage{
		Account: &AccountStore{rdb: rbd},
	}
}

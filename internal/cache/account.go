package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deriannavy/microgo/internal/store"
	"github.com/go-redis/redis/v8"
)

type AccountStore struct {
	rdb *redis.Client
}

const AccountExpTime = time.Minute

func (s *AccountStore) Get(ctx context.Context, accountId int64) (*store.Account, error) {
	cacheKey := fmt.Sprintf("account-%d", accountId)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var account store.Account
	if data != "" {
		err := json.Unmarshal([]byte(data), &account)
		if err != nil {
			return nil, err
		}
	}

	return &account, nil
}

func (s *AccountStore) Set(ctx context.Context, account *store.Account) error {
	cacheKey := fmt.Sprintf("account-%d", account.Id)

	json, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return s.rdb.SetEX(ctx, cacheKey, json, AccountExpTime).Err()
}

func (s *AccountStore) Delete(ctx context.Context, accountId int64) {
	cacheKey := fmt.Sprintf("account-%d", accountId)
	s.rdb.Del(ctx, cacheKey)
}

package repository

import (
	"time"

	"github.com/go-redis/redis"
)

type CacheRepo interface {
	Get(key string) (string, error)
	Set(key string, value string, expire time.Duration) error
	Refresh(key string, expiration time.Duration) error
}

type cacheRepo struct {
	redis *redis.Client
}

func NewCacheRepo(redis *redis.Client) CacheRepo {
	return &cacheRepo{
		redis: redis,
	}
}

func (r *cacheRepo) Get(key string) (string, error) {
	value, err := r.redis.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return value, nil
}

func (r *cacheRepo) Set(key string, value string, expire time.Duration) error {
	err := r.redis.Set(key, value, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *cacheRepo) Refresh(key string, expiration time.Duration) error {
	return r.redis.Expire(key, expiration).Err()
}

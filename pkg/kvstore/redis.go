package kvstore

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/nzwice/surl/pkg/config"
	"github.com/nzwice/surl/pkg/logging"
	"github.com/redis/go-redis/v9"
)

var (
	ErrRedisClient     = errors.New("redis client error")
	ErrDeSerialization = errors.New("serialization error")
	ErrEmptyValue      = errors.New("empty value")
)

type Client interface {
	DoGet(ctx context.Context, k string, v any) error
	SetTTL(ctx context.Context, k string, v any, ttl time.Duration) error
}

type redisImpl struct {
	client      *redis.Client
	serialize   func(v any) ([]byte, error)
	deserialize func(b []byte, v any) error
}

// DoGet implements [Client].
func (r *redisImpl) DoGet(ctx context.Context, k string, v any) error {
	ret, err := r.client.Get(ctx, k).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ErrEmptyValue
		}
		slog.ErrorContext(ctx, "redis client error", logging.ErrorAttr(err))
		return errors.Join(ErrRedisClient, err)
	}
	if err := r.deserialize([]byte(ret), v); err != nil {
		return errors.Join(ErrDeSerialization, err)
	}
	return nil
}

// SetTTL implements [Client].
func (r *redisImpl) SetTTL(ctx context.Context, k string, v any, ttl time.Duration) error {
	b, err := r.serialize(v)
	if err != nil {
		return errors.Join(ErrDeSerialization, err)
	}
	_, err = r.client.Set(ctx, k, string(b), ttl).Result()
	if err != nil {
		return errors.Join(ErrRedisClient, err)
	}
	return nil
}

func NewRedis(cfg config.RedisConfig) Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		DB:       cfg.DB,
		Password: cfg.Password,
	})
	return &redisImpl{
		client:      client,
		serialize:   json.Marshal,
		deserialize: json.Unmarshal,
	}
}

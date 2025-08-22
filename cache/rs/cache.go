package rs

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const RedisNil = redis.Nil

type ICache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value any, ex time.Duration) *redis.StatusCmd

	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	MSet(ctx context.Context, values ...interface{}) *redis.StatusCmd

	Pipeline() redis.Pipeliner
}

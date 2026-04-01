package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xpfo-go/logs"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
)

var DefaultRedisClient *redis.Client

func InitRedis(cfg *config.RedisConfig) {
	if cfg == nil {
		panic("redis config is nil")
	}
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		panic(fmt.Errorf("redis ping failed for %s: %w", addr, err))
	}

	DefaultRedisClient = client
	logs.Infof("connect to redis: %s", addr)
}

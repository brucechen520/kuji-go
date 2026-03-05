package pkg

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/redis/go-redis/v9"
)

// InitRedis 負責建立 Redis 連線
func InitRedis(cfg *config.Config) (*redis.Client, func(), error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	// 測試連線是否成功 (Fail Fast 原則)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	cleanup := func() {
		log.Println("正在關閉 Redis 連線...")
		if err := rdb.Close(); err != nil {
			log.Printf("Redis 關閉失敗: %v", err)
		} else {
			log.Println("Redis 已安全關閉")
		}
	}

	return rdb, cleanup, nil
}

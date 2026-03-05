package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/brucechen520/kuji-go/internal/repository/postgre"
	REDIS_REPO "github.com/brucechen520/kuji-go/internal/repository/redis"
)

// ProviderSet 是一張「零件清單」，告訴 Wire 這層樓有哪些工具可以用
var ProviderSet = wire.NewSet(
	InitDB,
	InitRedis,
	postgre.NewUserRepo,      // 這裡是你具體的 MySQL/Postgres 實作
	REDIS_REPO.NewTokenStore, // 這裡是你具體的 Redis 實作
)

// InitDB 負責建立 PostgreSQL 連線
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// 這裡的參數會從 config 讀取部署時注入的環境變數
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Taipei",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		cfg.DB.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 設定 Connection Pool (這也是重要的實作細節)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// InitRedis 負責建立 Redis 連線
func InitRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	// 測試連線是否成功 (Fail Fast 原則)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}

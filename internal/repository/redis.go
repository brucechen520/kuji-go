package repository // 屬於 repository 套件

import (
	"context" // 引入 context 套件，用於控制請求的生命週期 (如超時、取消)
	"fmt"     // 用於格式化字串
	"os"      // 用於讀取環境變數

	"github.com/redis/go-redis/v9" // 引入 Redis 客戶端
)

// NewRedis 初始化 Redis 客戶端
func NewRedis() (*redis.Client, error) {
	// 組合 Redis 位址 (Host:Port)
	addr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))

	// 建立 Redis 客戶端實例
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,                        // Redis 伺服器位址
		Password: os.Getenv("REDIS_PASSWORD"), // Redis 密碼 (若無則為空字串)
		DB:       0,                           // 使用預設的 DB 0
	})

	// 使用 Ping 檢查連線是否成功
	// context.Background() 產生一個空的 Context，表示沒有超時限制
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("無法連線至 Redis: %w", err) // %w 用於包裝錯誤 (Wrap Error)，保留原始錯誤訊息以便後續追蹤
	}
	return rdb, nil // 回傳成功的客戶端實例
}

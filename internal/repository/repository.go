package repository // 定義套件名稱為 repository，通常與資料夾名稱一致，負責資料存取層

import (
	"github.com/redis/go-redis/v9" // 引入 Redis 套件"
	"gorm.io/gorm"                 // 引入 GORM ORM 核心套件
)

// Repository 封裝資料庫操作，避免 Handler 直接依賴 GORM
// 定義一個結構體 (struct)，用來封裝資料庫連線
type Repository struct {
	db  *gorm.DB      // db 欄位儲存 GORM 的資料庫連線指標 (*gorm.DB)
	rdb *redis.Client // rdb 欄位儲存 Redis 的客戶端指標 (*redis.Client)
}

// NewRepository 是一個建構函式 (Constructor)，用於建立 Repository 實例
// 透過依賴注入 (Dependency Injection) 的方式傳入 db 連線
func NewRepository(db *gorm.DB, rdb *redis.Client) *Repository {
	return &Repository{db: db, rdb: rdb} // 回傳 Repository 結構體的指標，並初始化 db 與 rdb 欄位
}

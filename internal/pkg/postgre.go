package pkg

import (
	"fmt"
	"log"
	"time"

	"github.com/brucechen520/kuji-go/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 負責建立 PostgreSQL 連線
func InitDB(cfg *config.Config) (*gorm.DB, func(), error) {
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
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 設定 Connection Pool (這也是重要的實作細節)
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 定義具備日誌功能的 cleanup
	cleanup := func() {
		log.Println("正在關閉 PostgreSQL 連線...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("PostgreSQL 關閉失敗: %v", err)
		} else {
			log.Println("PostgreSQL 已安全關閉")
		}
	}

	return db, cleanup, nil
}

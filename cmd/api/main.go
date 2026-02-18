package main

import (
	"kuji-go/internal/models"
	"kuji-go/internal/repository"
	"kuji-go/internal/router"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// 1. 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("未發現 .env 檔案，使用系統環境變數")
	}

	// 2. 初始化資料庫 (PostgreSQL)
	repository.InitDB()
	// 自動遷移 Schema
	repository.DB.AutoMigrate(&models.Series{}, &models.Box{}, &models.Prize{})

	// 3. 初始化快取 (Redis)
	repository.InitRedis()

	// 4. 設定 Gin 路由
	r := router.SetupRouter()

	// 5. 啟動伺服器
	log.Println("一番賞系統成功啟動於 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("伺服器啟動失敗: ", err)
	}
}

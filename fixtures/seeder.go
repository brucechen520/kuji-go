package main

import (
	"kuji-go/internal/models"
	"kuji-go/internal/pkg"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// 載入環境變數，從專案根目錄的 .env 檔案
	// 假設 seeder 是從專案根目錄執行的 (e.g., go run internal/fixtures/seeder.go)
	if err := godotenv.Load(); err != nil {
		log.Println("找不到 .env 檔案，將使用系統環境變數")
	}

	// 建立資料庫連線
	db, err := pkg.NewDB()
	if err != nil {
		log.Fatalf("資料庫連線失敗: %v", err)
	}
	log.Println("資料庫連線成功！")

	// 執行資料遷移，確保資料表存在
	err = db.AutoMigrate(models.AllModels...)
	if err != nil {
		log.Fatalf("資料庫遷移失敗: %v", err)
	}
	log.Println("資料庫遷移成功！")

	// 清空資料表
	if err := clearTables(db); err != nil {
		log.Fatalf("清空資料表失敗: %v", err)
	}

	// 注入資料
	if err := seedData(db); err != nil {
		log.Fatalf("注入資料失敗: %v", err)
	}

	log.Println("✅ 資料注入成功！")
}

// clearTables 清空相關資料表
func clearTables(db *gorm.DB) error {
	log.Println("正在清空資料表...")
	// 依賴順序反向刪除，先刪 Prize，再刪 Box，最後刪 Series
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Prize{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Box{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.Series{}).Error; err != nil {
		return err
	}
	if err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&models.ProbabilityPhase{}).Error; err != nil {
		return err
	}
	log.Println("資料表清空完畢。")
	return nil
}

// seedData 讀取 JSON 檔案並寫入資料庫
func seedData(db *gorm.DB) error {
	data := GetSeries()

	// GORM 會自動遞迴處理：
	// Series -> Boxes -> Prizes -> ProbabilityPhases
	// 並且會自動將產生的 ID 填入對應的 SeriesID, BoxID, PrizeID
	return db.Create(&data).Error
}

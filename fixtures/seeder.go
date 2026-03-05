package main

import (
	"fmt"
	"log"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/brucechen520/kuji-go/internal/model"
	"github.com/brucechen520/kuji-go/internal/pkg"

	"gorm.io/gorm"
)

func main() {
	// 建立資料庫連線
	cfg, err := config.Load()

	if err != nil {
		fmt.Println("Load config file failed")
	}

	db, cleanup, err := pkg.InitDB(cfg)
	if err != nil {
		log.Fatalf("資料庫連線失敗: %v", err)
	}

	defer cleanup()

	log.Println("資料庫連線成功！")

	// 執行資料遷移，確保資料表存在
	err = db.AutoMigrate(model.AllModels...)
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
	log.Println("正在清空所有資料表並重置 ID...")

	// 這裡列出所有你想清空的資料表名稱（通常是複數，請根據你 DB 中的實際名稱調整）
	tables := []string{
		"probability_phases",
		"prizes",
		"boxes",
		"series",
		"users",
		"wallet_histories",
		"draw_logs",
	}

	for _, table := range tables {
		// 使用 TRUNCATE TABLE ... IF EXISTS ... CASCADE
		// RESTART IDENTITY 會讓自增 ID 回到 1
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s", table)).Error
		if err != nil {
			log.Printf("清空資料表 %s 失敗: %v", table, err)
			return err
		}
	}

	log.Println("✅ 資料表清空與 ID 重置完畢。")
	return nil
}

func seedData(db *gorm.DB) error {
	series := GetSeries()

	// GORM 會自動遞迴處理：
	// Series -> Boxes -> Prizes -> ProbabilityPhases
	// 並且會自動將產生的 ID 填入對應的 SeriesID, BoxID, PrizeID
	err := db.Create(&series).Error

	users := GetUsers()

	err = db.Create(&users).Error

	return err
}

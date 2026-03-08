package main

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/brucechen520/kuji-go/internal/model"
	"github.com/brucechen520/kuji-go/internal/pkg"

	"gorm.io/gorm"
)

func main() {
	seederLogger, _ := zap.NewProduction()
	defer seederLogger.Sync()

	// 建立資料庫連線
	cfg, err := config.Load()

	if err != nil {
		seederLogger.Fatal("Load config file failed", zap.Error(err))
	}

	db, cleanup, err := pkg.InitDB(cfg, seederLogger)
	if err != nil {
		seederLogger.Fatal("資料庫連線失敗", zap.Error(err))
	}

	defer cleanup()

	seederLogger.Info("資料庫連線成功！")

	// 執行資料遷移，確保資料表存在
	err = db.AutoMigrate(model.AllModels...)
	if err != nil {
		seederLogger.Fatal("資料庫遷移失敗", zap.Error(err))
	}
	seederLogger.Info("資料庫遷移成功！")

	// 清空資料表
	if err := clearTables(db, seederLogger); err != nil {
		seederLogger.Fatal("清空資料表失敗", zap.Error(err))
	}

	// 注入資料
	if err := seedData(db); err != nil {
		seederLogger.Fatal("注入資料失敗", zap.Error(err))
	}

	seederLogger.Info("✅ 資料注入成功！")
}

// clearTables 清空相關資料表
func clearTables(db *gorm.DB, seederLogger *zap.Logger) error {
	seederLogger.Info("正在清空所有資料表並重置 ID...")

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
		err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)).Error
		if err != nil {
			seederLogger.Error("清空資料表失敗", zap.String("table", table), zap.Error(err))
			return err
		}
	}

	seederLogger.Info("✅ 資料表清空與 ID 重置完畢。")
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

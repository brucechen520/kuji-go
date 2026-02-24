package app // 定義 app 套件，負責應用程式的啟動與依賴組裝

import (
	"kuji-go/internal/handlers"   // 引入 handlers 層
	"kuji-go/internal/models"     // 引入 models 層
	"kuji-go/internal/pkg"        // 引入 pkg 層
	"kuji-go/internal/repository" // 引入 repository 層
	"kuji-go/internal/service"    // 引入 service 層
	"log"                         // 引入日誌套件

	"github.com/joho/godotenv" // 引入環境變數載入工具
)

// Container 結構體，作為依賴注入容器 (DI Container)
// 它持有應用程式啟動後的所有核心元件
type Container struct {
	Handler *handlers.Handler // 總 Handler，包含所有子 Handler
}

// NewContainer 負責初始化所有元件並進行組裝
// 回傳 Container 以及一個 cleanup 函式，用於優雅關閉資源
func NewContainer() (*Container, func()) {
	// 1. 載入環境變數
	if err := godotenv.Load(); err != nil {
		log.Println("未發現 .env 檔案，使用系統環境變數") // 若無 .env 檔則忽略，繼續執行
	}

	// 2. 初始化資料庫連線
	db, err := pkg.NewDB() // 呼叫 pkg 層建立 DB 連線
	if err != nil {
		log.Fatal("資料庫連線失敗: ", err) // 若連線失敗，直接終止程式
	}
	// 自動遷移資料表結構 (Auto Migration)
	db.AutoMigrate(&models.Series{}, &models.Box{}, &models.Prize{})

	// 3. 初始化 Redis 連線
	rdb, err := pkg.NewRedis() // 呼叫 pkg 層建立 Redis 連線
	if err != nil {
		log.Fatal("Redis 連線失敗: ", err) // 若連線失敗，直接終止程式
	}

	// 4. 由下而上組裝依賴 (Wiring)

	// [Repository Layer] 建立 Repository
	repo := repository.NewRepository(db) // 注入 DB 連線

	// [Service Layer] 建立 Service，注入 Repository 和 Redis
	prizeService := service.NewPrizeService(repo, rdb) // 注入 repo 和 rdb

	// [Handler Layer] 建立 Handler，注入 Service
	prizeHandler := handlers.NewPrizeHandler(prizeService) // 注入 service

	// [Root Handler] 建立總 Handler，聚合所有子 Handler
	// 使用 Functional Options Pattern 進行組裝，更具可讀性與擴展性
	rootHandler := handlers.NewHandler(
		handlers.WithPrizeHandler(prizeHandler),
	)

	// 定義資源釋放邏輯
	cleanup := func() {
		log.Println("正在關閉應用程式資源...")

		// 關閉 Redis 連線
		pkg.CloseRedis(rdb)

		// 關閉 DB (GORM 需要先取得底層 sql.DB 才能關閉)
		if err := pkg.CloseDB(db); err != nil {
			log.Println("關閉資料庫失敗:", err)
		}
	}

	// 回傳包含完整依賴的容器
	return &Container{
		Handler: rootHandler, // 設定 Handler 欄位
	}, cleanup
}

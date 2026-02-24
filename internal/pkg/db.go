package pkg // 定義套件名稱為 pkg，通常與資料夾名稱一致，負責資料存取層

import (
	"fmt" // 引入 fmt 套件，用於字串格式化 (Format)
	"log"
	"os" // 引入 os 套件，用於操作作業系統功能，例如讀取環境變數
	"strconv"
	"time"

	"gorm.io/driver/postgres" // 引入 GORM 的 PostgreSQL 驅動程式
	"gorm.io/gorm"            // 引入 GORM ORM 核心套件
)

// NewDB 負責建立與資料庫的實體連線
// 回傳值為 (*gorm.DB, error)，這是 Go 的慣用寫法，同時回傳結果與錯誤
func NewDB() (*gorm.DB, error) {
	// 建議從環境變數讀取
	host := os.Getenv("DB_HOST")         // 讀取環境變數中的 DB_HOST (主機位址)
	user := os.Getenv("DB_USER")         // 讀取 DB_USER (使用者名稱)
	password := os.Getenv("DB_PASSWORD") // 讀取 DB_PASSWORD (密碼)
	dbname := os.Getenv("DB_NAME")       // 讀取 DB_NAME (資料庫名稱)
	port := os.Getenv("DB_PORT")         // 讀取 DB_PORT (連接埠)

	// 使用 fmt.Sprintf 格式化字串，組合成 PostgreSQL 的連線字串 (DSN - Data Source Name)
	// sslmode=disable 表示開發環境不使用 SSL 加密連線
	// TimeZone=Asia/Taipei 設定時區，確保時間存取正確
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		host, user, password, dbname, port)

	// gorm.Open 嘗試建立資料庫連線
	// postgres.Open(dsn) 指定使用 Postgres 驅動
	// &gorm.Config{} 可以傳入額外的 GORM 設定，這裡使用預設值
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Go 語言標準的錯誤處理方式：檢查 err 是否不為 nil
	if err != nil {
		return nil, err // 如果有錯誤，回傳 nil 連線物件和錯誤訊息
	}

	// --- 連線池 (Connection Pool) 設定 ---
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("無法取得底層 sql.DB: %w", err)
	}

	// 從環境變數讀取連線池設定，若無則使用預設值
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if maxIdleConns == 0 {
		maxIdleConns = 10 // 設定最大閒置連線數 (建議值)
	}

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if maxOpenConns == 0 {
		maxOpenConns = 100 // 設定最大開啟連線數 (建議值)
	}

	// 設定連線可被重複使用的最大時間
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour) // 建議設定為一小時，避免連線因網路問題失效

	return db, nil // 如果成功，回傳 db 連線物件和 nil 錯誤
}

func CloseDB(db *gorm.DB) error {
	// 2. 取得底層 sql.DB 物件
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("無法取得底層資料庫連線:", err)
		return err
	}

	// 3. 關閉連線並回傳結果
	log.Println("Closing database connection pool...")
	return sqlDB.Close()
}

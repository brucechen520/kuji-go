package repository // 定義套件名稱為 repository，通常與資料夾名稱一致，負責資料存取層

import (
	"fmt" // 引入 fmt 套件，用於字串格式化 (Format)
	"os"  // 引入 os 套件，用於操作作業系統功能，例如讀取環境變數

	"gorm.io/driver/postgres" // 引入 GORM 的 PostgreSQL 驅動程式
	"gorm.io/gorm"            // 引入 GORM ORM 核心套件
)

// Repository 封裝資料庫操作，避免 Handler 直接依賴 GORM
// 定義一個結構體 (struct)，用來封裝資料庫連線
type Repository struct {
	db *gorm.DB // db 欄位儲存 GORM 的資料庫連線指標 (*gorm.DB)
}

// NewRepository 是一個建構函式 (Constructor)，用於建立 Repository 實例
// 透過依賴注入 (Dependency Injection) 的方式傳入 db 連線
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db} // 回傳 Repository 結構體的指標，並初始化 db 欄位
}

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

	return db, nil // 如果成功，回傳 db 連線物件和 nil 錯誤
}

// WithTransaction 範例：封裝交易邏輯，確保 Lock 安全
// 這是一個高階函式，接收一個函式 fn 作為參數
func (r *Repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	// 呼叫 GORM 的 Transaction 方法
	// 它會自動開啟交易，執行 fn，如果 fn 回傳錯誤則 Rollback，否則 Commit
	return r.db.Transaction(fn)
}

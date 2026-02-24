package repository // 定義套件名稱為 repository，通常與資料夾名稱一致，負責資料存取層

import (
	"gorm.io/gorm" // 引入 GORM ORM 核心套件
)

// Repository 封裝資料庫操作，避免 Handler 直接依賴 GORM
// 定義一個結構體 (struct)，用來封裝資料庫連線
type Repository struct {
	db *gorm.DB // db 欄位儲存 GORM 的資料庫連線指標 (*gorm.DB)
}

// NewRepository 是一個建構函式 (Constructor)，用於建立 Repository 實例
// 透過依賴注入 (Dependency Injection) 的方式傳入 db 連線
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db} // 回傳 Repository 結構體的指標，並初始化 db 與 rdb 欄位
}

// WithTransaction 範例：封裝交易邏輯，確保 Lock 安全
// 這是一個高階函式，接收一個函式 fn 作為參數
func (r *Repository) WithTransaction(fn func(tx *gorm.DB) error) error {
	// 呼叫 GORM 的 Transaction 方法
	// 它會自動開啟交易，執行 fn，如果 fn 回傳錯誤則 Rollback，否則 Commit
	return r.db.Transaction(fn)
}

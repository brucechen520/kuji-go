package pkg

import "gorm.io/gorm"

// TransactionManager 封裝了資料庫事務操作。
// 它可以被注入到需要事務能力的服務中，而不需要傳遞整個 Repository。
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 建立一個新的事務管理器。
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{db: db}
}

// Exec 執行一個事務。
// 它會傳入一個 `*gorm.DB` 的事務物件給 `fn` 函式。
// 如果 `fn` 回傳任何錯誤，事務將會自動回滾 (Rollback)。
// 如果 `fn` 執行成功 (回傳 nil)，事務將會自動提交 (Commit)。
func (tm *TransactionManager) Exec(fn func(tx *gorm.DB) error) error {
	return tm.db.Transaction(fn)
}

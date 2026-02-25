package models

import "gorm.io/gorm"

// WalletHistory 錢包流水帳 (帳務審計用)
type WalletHistory struct {
	gorm.Model
	UserID      uint   `gorm:"index"`
	ActionType  string // "TOPUP"(儲值), "DRAW"(抽獎消費), "REFUND"(退款)
	Amount      int    // 變動金額 (正負值)
	Balance     int    // 變動後的餘額快照
	Description string // 備註
}

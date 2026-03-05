package model // 定義 models 套件，存放資料結構

import "gorm.io/gorm" // 引入 GORM 套件

type User struct {
	gorm.Model
	Username      string `gorm:"unique;not null"`
	Email         string `gorm:"unique"`
	Wallet        Wallet
	DrawLogs      []DrawLog       // 抽獎紀錄
	WalletHistory []WalletHistory // 點數增減流水帳
}

// DrawLog 抽獎紀錄 (中獎快照)
type DrawLog struct {
	gorm.Model
	UserID     uint   `gorm:"index"`
	BoxID      uint   `gorm:"index"`
	PrizeID    uint   `gorm:"index"`
	PaidPrice  int    `gorm:"not null"` // 紀錄抽獎當下花多少點數
	PrizeName  string `gorm:"not null"` // 紀錄中獎獎項名稱快照
	PrizeLevel string `gorm:"not null"` // 紀錄等級快照 (A, B...)
}

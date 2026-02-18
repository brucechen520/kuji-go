package models

import "gorm.io/gorm"

// 一番賞系列
type Series struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Boxes       []Box // 一對多關聯
}

// 抽獎箱
type Box struct {
	gorm.Model
	SeriesID       uint
	LocationName   string `gorm:"index"`
	TotalQuantity  int
	RemainQuantity int
	Prizes         []Prize
}

// 獎項
type Prize struct {
	gorm.Model
	BoxID      uint
	Level      string `gorm:"type:varchar(10)"` // A, B, C...
	Name       string
	Weight     int // 機率權重
	RemainStep int // 剩餘數量
}

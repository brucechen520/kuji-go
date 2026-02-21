package models // 定義 models 套件，存放資料結構

import "gorm.io/gorm" // 引入 GORM 套件

// 一番賞系列
// 定義 Series 結構體，對應資料庫中的 series 表
type Series struct {
	gorm.Model         // 內嵌 gorm.Model，自動包含 ID, CreatedAt, UpdatedAt, DeletedAt 欄位
	Name        string `gorm:"unique;not null"` // 設定 Name 欄位為唯一且不可為空
	Description string // 描述欄位
	Boxes       []Box  // 定義一對多關聯 (Has Many)，一個系列可以有多個箱子
}

// 抽獎箱
// 定義 Box 結構體，對應 boxes 表
type Box struct {
	gorm.Model             // 內嵌基礎欄位
	SeriesID       uint    // 外鍵 (Foreign Key)，關聯到 Series 表
	LocationName   string  `gorm:"index"` // 建立索引，加快查詢速度
	TotalQuantity  int     // 總數量
	RemainQuantity int     // 剩餘數量
	Prizes         []Prize // 定義一對多關聯，一箱有多個獎品
}

// 獎項
// 定義 Prize 結構體，對應 prizes 表
type Prize struct {
	gorm.Model        // 內嵌基礎欄位
	BoxID      uint   // 外鍵，關聯到 Box 表
	Level      string `gorm:"type:varchar(10)"` // 指定資料庫欄位型態為 varchar(10)，例如 "A", "LastOne"
	Name       string // 獎品名稱
	Weight     int    // 機率權重 (用於計算抽中機率)
	RemainStep int    // 該獎項剩餘數量
}

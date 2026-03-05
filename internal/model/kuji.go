package model // 定義 models 套件，存放資料結構

import "gorm.io/gorm" // 引入 GORM 套件

// 一番賞系列
// 定義 Series 結構體，對應資料庫中的 series 表
type Series struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
	Price       int   `gorm:"not null;default:0"` // 每一抽需要的代幣點數
	Boxes       []Box // 定義一對多關聯 (Has Many)，一個系列可以有多個箱子
}

// 抽獎箱
// 定義 Box 結構體，對應 boxes 表
// 已抽數可由 (TotalQuantity - RemainQuantity) 計算得出
type Box struct {
	gorm.Model
	SeriesID       uint
	LocationName   string    `gorm:"index"`
	TotalQuantity  int       // 總抽數，例如 80
	RemainQuantity int       // 剩餘抽數
	Prizes         []Prize   // Has Many 關聯
	DrawLogs       []DrawLog // 關聯：這箱產出的所有紀錄
}

// 獎項
// 定義 Prize 結構體，對應 prizes 表
type Prize struct {
	gorm.Model
	BoxID             uint
	Level             string `gorm:"type:varchar(10)"` // 例如 "A", "B", "LastOne"
	Name              string
	InitialQuantity   int                // 此獎項的初始總數，例如 A賞 1 個, B賞 2 個
	RemainingQuantity int                // 此獎項的剩餘數量
	Phases            []ProbabilityPhase // Has Many 關聯，用於動態機率
}

// ProbabilityPhase 定義了獎品在不同抽數階段的機率權重
type ProbabilityPhase struct {
	gorm.Model
	PrizeID        uint
	StartDrawCount int `gorm:"not null"` // 當已抽出 N 張時，此階段開始 (包含)
	EndDrawCount   int `gorm:"not null"` // 當已抽出 N 張時，此階段結束 (不包含)
	Weight         int `gorm:"not null"` // 在此階段的抽獎權重
}

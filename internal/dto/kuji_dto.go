package dto

import "time"

type SeriesDetailDTO struct {
	ID    uint           `json:"id"`
	Name  string         `json:"name"`
	Price int            `json:"price"`
	Boxes []BoxDetailDTO `json:"boxes"`
}

// BoxBaseDTO: 用於簡單列表，只需顯示箱子基本資訊
type BoxBaseDTO struct {
	ID           uint   `json:"id"`
	LocationName string `json:"location_name"`
}

// BoxDetailDTO: 用於一番賞詳細頁，包含該箱子下的所有獎項
type BoxDetailDTO struct {
	ID           uint             `json:"id"`
	LocationName string           `json:"location_name"`
	Prizes       []PrizeDetailDTO `json:"prizes"`
}

// 基礎屬性 (所有地方都要用)
type PrizeBaseDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// 詳細頁專用 (包含庫存)
type PrizeDetailDTO struct {
	ID                uint   `json:"id"`
	Name              string `json:"name"`
	Level             string `json:"level"`
	RemainingQuantity int    `json:"remaining_quantity"`
}

// 紀錄頁專用 (不包含庫存，但可能包含抽獎時間)
type PrizeHistoryDTO struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	DrawTime time.Time `json:"draw_time"`
}

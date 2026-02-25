package main

import "kuji-go/internal/models"

// PrizeNames 定義一箱中各賞別的名稱
type PrizeNames struct {
	A, B, C, D, E, LastOne string
}

// GetUsers 產生測試使用者，包含不同餘額情境
func GetUsers() []models.User {
	return []models.User{
		{
			Username: "admin",
			Email:    "admin@kuji.go",
			Balance:  99999, // 課長等級：測試大量抽獎
		},
		{
			Username: "lucky_guy",
			Email:    "lucky@example.com",
			Balance:  1000, // 普通玩家：測試正常抽獎扣款
		},
		{
			Username: "poor_student",
			Email:    "poor@example.com",
			Balance:  10, // 窮學生：測試餘額不足報錯
		},
	}
}

// GetSeries 產生系列資料，並賦予每抽價格
func GetSeries() []models.Series {
	OnePieceNames := PrizeNames{
		A: "魯夫模型", B: "索隆模型", C: "娜美大毛巾",
		D: "西門限定吊飾", E: "隨機色紙", LastOne: "凱多特別版",
	}

	DragonBallNames := PrizeNames{
		A: "自在極意功悟空", B: "深藍貝吉達", C: "悟吉塔毛巾",
		D: "戰鬥力偵測器吊飾", E: "七龍珠插畫板", LastOne: "神龍模型",
	}

	return []models.Series{
		{
			Name:        "海賊王-新世界篇",
			Description: "全台兩店同步開抽！",
			Price:       250, // 每一抽 250 代幣
			Boxes: []models.Box{
				generateStandardBox("台北西門店", OnePieceNames),
				generateStandardBox("台中逢甲店", OnePieceNames),
			},
		},
		{
			Name:        "七龍珠-超極限戰鬥",
			Description: "高雄限定場次",
			Price:       300, // 每一抽 300 代幣
			Boxes: []models.Box{
				generateStandardBox("高雄夢時代店", DragonBallNames),
			},
		},
	}
}

// generateStandardBox 輔助函式，自動配置 80 抽的獎項分佈與機率
func generateStandardBox(location string, names PrizeNames) models.Box {
	return models.Box{
		LocationName:   location,
		TotalQuantity:  80,
		RemainQuantity: 80,
		Prizes: []models.Prize{
			{
				Level: "A", Name: names.A, InitialQuantity: 1, RemainingQuantity: 1,
				Phases: []models.ProbabilityPhase{
					{StartDrawCount: 0, EndDrawCount: 20, Weight: 0},   // 護航機制：前 20 抽不會出 A
					{StartDrawCount: 20, EndDrawCount: 80, Weight: 10}, // 20 抽後權重開啟
				},
			},
			{
				Level: "B", Name: names.B, InitialQuantity: 2, RemainingQuantity: 2,
				Phases: []models.ProbabilityPhase{{StartDrawCount: 0, EndDrawCount: 80, Weight: 20}},
			},
			{
				Level: "C", Name: names.C, InitialQuantity: 10, RemainingQuantity: 10,
				Phases: []models.ProbabilityPhase{{StartDrawCount: 0, EndDrawCount: 80, Weight: 100}},
			},
			{
				Level: "D", Name: names.D, InitialQuantity: 27, RemainingQuantity: 27,
				Phases: []models.ProbabilityPhase{{StartDrawCount: 0, EndDrawCount: 80, Weight: 270}},
			},
			{
				Level: "E", Name: names.E, InitialQuantity: 40, RemainingQuantity: 40,
				Phases: []models.ProbabilityPhase{{StartDrawCount: 0, EndDrawCount: 80, Weight: 400}},
			},
			{
				Level: "LastOne", Name: names.LastOne, InitialQuantity: 1, RemainingQuantity: 1,
				Phases: []models.ProbabilityPhase{{StartDrawCount: 0, EndDrawCount: 80, Weight: 0}}, // 不參與機率分配
			},
		},
	}
}

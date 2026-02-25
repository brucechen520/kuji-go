package main

import "kuji-go/internal/models"

// PrizeNames 定義一箱中各賞別的名稱
type PrizeNames struct {
	A, B, C, D, E, LastOne string
}

// 在 data.go 中新增
func GetUsers() []models.User {
	return []models.User{
		{
			Username: "admin",
			Email:    "admin@kuji.go",
			Balance:  99999, // 測試用大戶
		},
		{
			Username: "lucky_guy",
			Email:    "lucky@example.com",
			Balance:  1000,
		},
		{
			Username: "poor_student",
			Email:    "poor@example.com",
			Balance:  10, // 測試餘額不足
		},
	}
}

func GetSeries() []models.Series {
	return []models.Series{
		{
			Name:        "海賊王-新世界篇",
			Description: "全台兩店同步開抽！",
			Boxes: []models.Box{
				generateStandardBox("台北西門店", PrizeNames{
					A: "魯夫模型", B: "索隆模型", C: "娜美大毛巾",
					D: "西門限定吊飾", E: "隨機色紙", LastOne: "凱多特別版",
				}),
				generateStandardBox("台中逢甲店", PrizeNames{
					A: "魯夫模型", B: "香吉士模型", C: "羅賓大毛巾",
					D: "逢甲限定胸章", E: "隨機色紙", LastOne: "大媽特別版",
				}),
			},
		},
		{
			Name: "七龍珠-超極限戰鬥",
			Boxes: []models.Box{
				generateStandardBox("高雄夢時代店", PrizeNames{
					A: "自在極意功悟空", B: "深藍貝吉達", C: "悟吉塔毛巾",
					D: "戰鬥力偵測器吊飾", E: "七龍珠插畫板", LastOne: "神龍模型",
				}),
			},
		},
	}
}

// 修改後的輔助函式，支援名稱自定義
func generateStandardBox(location string, names PrizeNames) models.Box {
	return models.Box{
		LocationName:   location,
		TotalQuantity:  80,
		RemainQuantity: 80,
		Prizes: []models.Prize{
			{
				Level: "A", Name: names.A, InitialQuantity: 1, RemainingQuantity: 1,
				Phases: []models.ProbabilityPhase{
					{StartDrawCount: 0, EndDrawCount: 20, Weight: 0}, // 護航
					{StartDrawCount: 20, EndDrawCount: 80, Weight: 10},
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
				Phases: []models.ProbabilityPhase{{StartDrawCount: 0, EndDrawCount: 80, Weight: 0}},
			},
		},
	}
}

package pkg

import "fmt"

// 未來如果有 by user需求 或 global需求的 cache的話，這裡要擴充成 folder

const (
	// 基礎前綴，方便未來做多租戶或環境隔離
	prefix  = "kuji"
	version = "v1"

	// 格式範本
	// 獎項清單快取: kuji:v1:series:{series_id}:prizes
	fmtSeriesPrizes = "series:%d:prizes"

	// 實時庫存計數: kuji:v1:box:{box_id}:inventory
	fmtBoxInventory = "box:%d:inventory"

	// 一番賞 Meta: kuji:v1:series:{series_id}:meta
	fmtSeriesMeta = "series:%d:meta"
)

// 提供輔助函數來生成具體的 Key
func GetSeriesPrizesKey(seriesID uint) string {
	return fmt.Sprintf("%s:%s:%s", prefix, version, fmt.Sprintf(fmtSeriesPrizes, seriesID))
}

func GetBoxInventoryKey(boxID uint) string {
	return fmt.Sprintf("%s:%s:%s", prefix, version, fmt.Sprintf(fmtBoxInventory, boxID))
}

func GetSeriesMetaKey(seriesID uint) string {
	return fmt.Sprintf("%s:%s:%s", prefix, version, fmt.Sprintf(fmtSeriesMeta, seriesID))
}

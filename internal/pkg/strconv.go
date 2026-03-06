package pkg

import "strconv"

// 這才是 Go 開發者常見的封裝邏輯
func ParseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0 // 或者 log.Printf("轉型錯誤: %v", err)
	}
	return val
}

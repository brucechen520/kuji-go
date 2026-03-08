package pkg

import "strconv"

// ParseInt 把 string 轉換成 int
func ParseInt(str string) int {
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0 // 或者回傳錯誤
	}
	return res
}

// StringToUint 把 string 轉換成 uint
// 如果有錯誤就預設回傳 0 (這是一種選擇，也可以選擇回傳 error)
func StringToUint(str string) uint {
	res, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0 // 或者回傳錯誤
	}
	return uint(res)
}

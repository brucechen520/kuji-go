// repository/repository.go
package repository

import (
	"errors"
)

// 定義 Repository 層級的通用錯誤
var (
	ErrRecordNotFound = errors.New("database record not found")
	ErrUpdateFailed   = errors.New("database update failed")
)

// 如果你有通用的分頁請求，也可以放這
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

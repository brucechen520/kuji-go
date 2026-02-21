package repository // 屬於 repository 套件

import (
	"context"                 // 用於傳遞 Context
	"kuji-go/internal/models" // 引入資料模型定義
)

// GetPrizesByBoxID 屬於 Prize 相關的 DB 操作
// (r *Repository) 表示這是 Repository 結構體的方法 (Method Receiver)
func (r *Repository) GetPrizesByBoxID(ctx context.Context, boxID string) ([]models.Prize, error) {
	var prizes []models.Prize // 宣告一個 slice 來存放查詢結果

	// r.db.WithContext(ctx): 將 Gin 的 Context 傳入 GORM，這樣如果 HTTP 請求被取消，DB 查詢也會被中斷
	// .Where("box_id = ?", boxID): SQL 的 WHERE 條件，使用 ? 防止 SQL Injection
	// .Find(&prizes): 執行查詢並將結果映射 (Map) 到 prizes 變數中
	result := r.db.WithContext(ctx).Where("box_id = ?", boxID).Find(&prizes)
	return prizes, result.Error // 回傳查詢到的資料與可能的錯誤
}

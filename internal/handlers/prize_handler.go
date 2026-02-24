package handlers // 定義套件名稱為 handlers，負責處理 HTTP 請求邏輯

import (
	"context"  // 引入 service 套件，使用商業邏輯
	"log"      // 引入日誌套件
	"net/http" // 引入標準庫 net/http

	"github.com/gin-gonic/gin" // 引入 Gin Web Framework
)

// PrizeService 是一個介面，定義了 PrizeHandler 所需的商業邏輯方法。
// 這樣 PrizeHandler 就只依賴於抽象，而不是具體的 service 實作。
type PrizeService interface {
	GetPrizes(ctx context.Context, boxID string) ([]string, error)
	Draw(ctx context.Context) (string, error)
}

// PrizeHandler 結構體負責聚合所有獎品相關的依賴
type PrizeHandler struct {
	Service PrizeService // 改為依賴介面
}

// NewPrizeHandler 是建構函式，用於初始化 PrizeHandler
// 透過參數傳入依賴 (Dependency Injection)，方便測試與管理
func NewPrizeHandler(s PrizeService) *PrizeHandler {
	return &PrizeHandler{Service: s} // 注入 Service 實例
}

// GetList 是一個方法 (Method)，綁定在 PrizeHandler 上
// 參數 c *gin.Context 包含了該次 HTTP 請求的所有資訊 (參數、標頭等)
func (h *PrizeHandler) GetList(c *gin.Context) {
	boxID := c.Query("box_id") // 從 URL Query String 獲取參數 (例如: /prizes?box_id=123)

	// 呼叫 Service 層的方法來獲取資料
	prizes, err := h.Service.GetPrizes(c.Request.Context(), boxID)
	if err != nil {
		// 在伺服器端記錄詳細錯誤，以便追蹤問題
		log.Printf("ERROR: h.Service.GetPrizes failed: %v\n", err)
		// 回傳給客戶端一個通用的錯誤，避免洩漏內部實作細節
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法獲取獎品列表"})
		return
	}

	// c.JSON 回傳 JSON 格式的回應
	// http.StatusOK 代表 HTTP 200
	// gin.H 是一個 map[string]interface{} 的捷徑，用來建構 JSON 物件
	c.JSON(http.StatusOK, gin.H{
		"box_id": boxID,  // 回傳查詢的 box_id
		"prizes": prizes, // 回傳從 Service 取得的獎品列表
	})
}

// Draw 處理抽獎請求
func (h *PrizeHandler) Draw(c *gin.Context) {
	// 呼叫 Service 層的抽獎方法
	msg, err := h.Service.Draw(c.Request.Context())
	if err != nil {
		log.Printf("ERROR: h.Service.Draw failed: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "系統繁忙，請稍後再試"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": msg, // 回傳 Service 處理後的結果
	})
}

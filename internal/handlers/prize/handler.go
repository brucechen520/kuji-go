package prize

import (
	"net/http"

	"github.com/brucechen520/kuji-go/internal/code"
	"github.com/brucechen520/kuji-go/internal/models"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"github.com/brucechen520/kuji-go/internal/service/client/prize"
	"go.uber.org/zap"
)

type Handler struct {
	service prize.Service
	logger  *zap.Logger
}

func New(s prize.Service, logger *zap.Logger) *Handler {
	return &Handler{service: s, logger: logger}
}

// ListPrizes 取得獎項列表的 Handler 實作
func (h *Handler) ListPrizes(ctx core.Context) {
	// 1. 驗證與綁定參數 (Request Binding)
	// 假設未來會有分頁或分類篩選，這裡先定義一個空的 Request 結構
	type ListRequest struct {
		Category string `form:"category"` // 從 Query String 取得 ?category=xxx
		Page     int    `form:"page,default=1"`
	}

	req := new(ListRequest)
	if err := ctx.ShouldBindQuery(req); err != nil {
		// 如果參數解析失敗，直接中斷並回傳錯誤 (由 core.go 的 middleware 處理)
		ctx.AbortWithError(core.Error(
			http.StatusBadRequest,
			code.HashIdsEncodeError,
			code.Text(code.HashIdsEncodeError)).WithError(err),
		)
		return
	}

	// 2. 呼叫 Service (目前先回傳固定 Slice)
	// 在真實情境中，你會呼叫: prizes, err := h.service.GetPrizeList()

	prizes := []models.Prize{
		{
			BoxID:             1,
			Name:              "一番賞 A 賞 - 孫悟空模型",
			Level:             "A",
			InitialQuantity:   1,
			RemainingQuantity: 1,
		},
		{
			BoxID:             2,
			Name:              "一番賞 B 賞 - 貝吉塔模型",
			Level:             "B",
			InitialQuantity:   2,
			RemainingQuantity: 2,
		},
		{
			BoxID:             3,
			Name:              "一番賞 C 賞 - 紀念毛巾",
			Level:             "C",
			InitialQuantity:   50,
			RemainingQuantity: 50,
		},
	}

	// 3. 回傳 Payload
	// ctx.Payload 會將資料存入自定義 Context 的 Payload 中，
	// 最後由 core.go 裡的 Middleware 統一執行 ctx.JSON(200, response)
	ctx.Payload(prizes)
}

package client

import (
	"net/http"

	clientS "github.com/brucechen520/kuji-go/internal/service/client"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"github.com/gin-gonic/gin"
)

type SeriesHandler struct {
	seriesService clientS.SeriesService
}

func NewSeriesHandler(as clientS.SeriesService) *SeriesHandler {
	return &SeriesHandler{seriesService: as}
}

func (s *SeriesHandler) GetSeriesById(c *gin.Context) {
	// 1. 解析參數 (這裡可以使用一個特定的 Request struct)
	var req struct {
		SeriesID uint `uri:"SeriesID" binding:"required,gt=0"`
	}
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 取出由 Middleware 建立好的 core.Context (帶有 trace_id 的 Logger)
	ctx := core.MustGetContext(c)

	result, err := s.seriesService.GetSeriesById(ctx, req.SeriesID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	// 3. 回傳結果
	c.JSON(http.StatusOK, gin.H{"series": result})
}

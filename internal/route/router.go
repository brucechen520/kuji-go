package route

import (
	"github.com/brucechen520/kuji-go/internal/config"
	clientH "github.com/brucechen520/kuji-go/internal/handler/client"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 也可以這樣設計
type handlerGroup struct {
	Auth   *clientH.AuthHandler
	Series *clientH.SeriesHandler
}

// 定義這個組裝函數
func NewHandlerGroup(auth *clientH.AuthHandler, series *clientH.SeriesHandler) *handlerGroup {
	return &handlerGroup{
		Auth:   auth,
		Series: series,
	}
}

// 註冊所有路由
func RegisterRoutes(e *gin.Engine, h *handlerGroup) {
	v1 := e.Group("/api/v1/client")
	{
		v1.POST("/login", h.Auth.Login)
		v1.GET("/series/:SeriesID/prizes", h.Series.GetSeriesById)
	}
}

// NewHTTPServer 直接在這裡組裝 Engine 與路由
func NewHTTPServer(logger *zap.Logger, cfg *config.AppConfig, h *handlerGroup) *core.Engine {
	// 1. 初始化引擎，注入 Config 選項
	engine := core.NewEngine(logger, cfg)

	// 2. 註冊路由 (直接操作 engine.Engine)
	RegisterRoutes(engine.Engine, h)

	return engine
}

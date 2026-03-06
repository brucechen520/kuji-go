package route

import (
	clientH "github.com/brucechen520/kuji-go/internal/handler/client"
	"github.com/gin-gonic/gin"
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

func NewRouter(h *handlerGroup) *gin.Engine {
	r := gin.Default()

	v1ClientGroup := r.Group("/api/v1/client")
	{
		v1ClientGroup.POST("/login", h.Auth.Login)
		v1ClientGroup.GET("/series/:SeriesID/prizes", h.Series.GetSeriesById)
	}

	return r
}

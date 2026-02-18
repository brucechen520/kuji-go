package router

import (
	"kuji-go/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 前台 API
	v1 := r.Group("/api/v1")
	{
		v1.GET("/prizes", handlers.GetPrizeList) // 獎項列表
		v1.POST("/draw", handlers.Draw)          // 抽獎
		v1.GET("/user/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"user": "小明", "coins": 100})
		})
	}

	// 後台管理 API (可以加 Middleware 做認證)
	admin := r.Group("/admin")
	{
		admin.POST("/box/init", func(c *gin.Context) {
			c.JSON(201, gin.H{"status": "海賊王箱子已初始化"})
		})
		admin.PATCH("/prize/probability", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "機率調整成功"})
		})
	}

	return r
}

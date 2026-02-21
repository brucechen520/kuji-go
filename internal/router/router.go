package router // 定義 router 套件

import (
	"kuji-go/internal/handlers" // 引入 handlers 套件

	"github.com/gin-gonic/gin" // 引入 Gin 框架
)

// SetupRouter 設定所有的路由規則
// 接收一個初始化好的 Handler 結構體
func SetupRouter(h *handlers.Handler) *gin.Engine {
	r := gin.Default() // 建立一個預設的 Gin 引擎 (包含 Logger 和 Recovery 中介軟體)

	// 前台 API
	v1 := r.Group("/api/v1") // 建立路由群組，所有底下的 API 都會以 /api/v1 開頭
	{
		// 定義 GET /api/v1/prizes，並指定處理函式為 h.Prize.GetList
		v1.GET("/prizes", h.Prize.GetList)

		// 定義 POST /api/v1/draw，處理抽獎
		v1.POST("/draw", h.Prize.Draw)

		// 定義一個簡單的測試路由，使用匿名函式 (Anonymous Function) 直接處理
		v1.GET("/user/me", func(c *gin.Context) {
			c.JSON(200, gin.H{"user": "小明", "coins": 100})
		})
	}

	// 後台管理 API (可以加 Middleware 做認證)
	admin := r.Group("/admin") // 建立後台路由群組 /admin
	{
		admin.POST("/box/init", func(c *gin.Context) {
			c.JSON(201, gin.H{"status": "海賊王箱子已初始化"})
		})
		admin.PATCH("/prize/probability", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "機率調整成功"})
		})
	}

	return r // 回傳設定好的 Gin 引擎
}

package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Context 定義我們自己的 Web 上下文規範
type Context interface {
	// 綁定資料
	ShouldBindJSON(obj interface{}) error

	// 統一回傳：成功與失敗
	Payload(payload interface{})
	AbortWithError(err error) // 這裡簡化，直接傳 error

	// 獲取原生的 Context (傳遞給 Service 用)
	RequestContext() gin.Context
}

type context struct {
	ctx *gin.Context
}

// 實作介面方法
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindJSON(obj)
}

func (c *context) Payload(payload interface{}) {
	// 第一階段：在這裡直接寫死統一的 JSON 結構
	c.ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "OK",
		"data": payload,
	})
}

func (c *context) AbortWithError(err error) {
	// 第一階段：簡單處理錯誤回傳
	c.ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"code": 500,
		"msg":  err.Error(),
	})
}

func (c *context) RequestContext() gin.Context {
	return *c.ctx
}

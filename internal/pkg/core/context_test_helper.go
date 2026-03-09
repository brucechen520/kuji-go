package core

import "github.com/gin-gonic/gin"

// NewTestContext 僅供測試使用，從 *gin.Context 建立一個完整的 core.Context。
// 生產環境中 Context 的建立和釋放由 LoggerMiddleware 統一管理（透過 contextPool）；
// 測試環境中因為沒有 Middleware，需要用此函式手動建立。
//
// 典型用法：
//
//	c, _ := gin.CreateTestContext(nil)
//	ctx := core.NewTestContext(c)
func NewTestContext(ctx *gin.Context) Context {
	c := contextPool.Get().(*customContext)
	c.ctx = ctx
	return c
}

// ReleaseTestContext 搭配 NewTestContext 使用，歸還 Pool 並清理。
// 建議配合 defer 呼叫：
//
//	ctx := core.NewTestContext(c)
//	defer core.ReleaseTestContext(ctx)
func ReleaseTestContext(c Context) {
	impl := c.(*customContext)
	impl.ctx = nil
	contextPool.Put(impl)
}

package core

import (
	"bytes"
	"context"
	"io"
	"sync"

	"github.com/brucechen520/kuji-go/pkg/trace"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	_LoggerName  = "_logger_"
	_TraceName   = "_trace_"
	_BodyName    = "_body_"
	_coreCtxKey  = "_core_ctx_" // middleware 建立好的 core.Context 存放位置
)

// Context 介面，僅暴露「讀取行為」給外部 (Service / Handler / Repository) 使用。
// setLogger / setTrace 等設定行為只由 core package 內部透過 *customContext 直接呼叫，不屬於介面的職責。
// 這樣做的好處是外部 package 可以自由實作或 Mock 此介面，不受未導出方法的限制。
type Context interface {
	GetLogger() *zap.Logger
	Trace() trace.T
	RawData() []byte
	StdContext() context.Context
}

type customContext struct {
	ctx *gin.Context
}

// contextPool 復用 *customContext 的記憶體，避免每次 request 都重新分配。
// 生命週期完全由 LoggerMiddleware 管理：請求開始時 Get()，結束時 Put()。
// 外部 package 不可直接存取，應透過 MustGetContext 取得已初始化的 Context。
var contextPool = &sync.Pool{
	New: func() interface{} {
		return &customContext{}
	},
}

// MustGetContext 從 gin.Context 中取出由 Middleware 已建立好的 core.Context。
// 使用此函式的前提是 core.LoggerMiddleware 已套用；若未套用則 panic。
// Handler 應一律使用此函式，而非自行呼叫 NewContext，以確保帶有 trace_id 的 Logger 正確傳遞。
func MustGetContext(c *gin.Context) Context {
	v, exists := c.Get(_coreCtxKey)
	if !exists {
		panic("core.Context not found: ensure core.LoggerMiddleware is applied before this handler")
	}
	return v.(Context)
}

func (c *customContext) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	c.ctx.Set(_BodyName, body)                               // cache body是為了trace使用
	c.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // re-construct req body
}

func (c *customContext) setLogger(l *zap.Logger) {
	c.ctx.Set(_LoggerName, l)
}

func (c *customContext) GetLogger() *zap.Logger {
	if l, ok := c.ctx.Get(_LoggerName); ok {
		return l.(*zap.Logger)
	}
	return zap.NewNop()
}

func (c *customContext) setTrace(t trace.T) {
	c.ctx.Set(_TraceName, t)
}

func (c *customContext) Trace() trace.T {
	if t, ok := c.ctx.Get(_TraceName); ok {
		return t.(trace.T)
	}
	return nil
}

func (c *customContext) RawData() []byte {
	body, ok := c.ctx.Get(_BodyName)
	if !ok {
		return nil
	}

	return body.([]byte)
}

func (c *customContext) StdContext() context.Context {
	if c.ctx == nil || c.ctx.Request == nil {
		return context.Background()
	}
	return c.ctx.Request.Context()
}

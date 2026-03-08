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
	_LoggerName = "_logger_"
	_TraceName  = "_trace_"
	_BodyName   = "_body_"
)

type Trace trace.T

// Context 介面，隱藏實作細節
type Context interface {
	setLogger(logger *zap.Logger)
	GetLogger() *zap.Logger
	setTrace(t Trace)
	Trace() Trace
	RawData() []byte
	StdContext() context.Context
}

type customContext struct {
	ctx *gin.Context
}

// sync.Pool 管理物件池
var contextPool = &sync.Pool{
	New: func() interface{} {
		return &customContext{}
	},
}

func NewContext(ctx *gin.Context) Context {
	c := contextPool.Get().(*customContext)
	c.ctx = ctx
	return c
}

func ReleaseContext(c Context) {
	impl := c.(*customContext)
	impl.ctx = nil
	contextPool.Put(impl)
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

func (c *customContext) setTrace(t Trace) {
	c.ctx.Set(_TraceName, t)
}

func (c *customContext) Trace() Trace {
	if t, ok := c.ctx.Get(_TraceName); ok {
		return t.(Trace)
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

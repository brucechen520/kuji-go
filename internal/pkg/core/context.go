package core

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"sync"

	stdctx "context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// HandlerFunc 定義了使用自定義 Context 的處理函式格式
type HandlerFunc func(c Context)

// 定義 Context 中使用的 Key，避免字串硬編碼（Hard-coding）
const (
	_Alias           = "_alias_"
	_TraceName       = "_trace_"
	_LoggerName      = "_logger_"
	_BodyName        = "_body_"
	_PayloadName     = "_payload_"
	_AbortErrorName  = "_abort_error_"
	_IsRecordMetrics = "_is_record_metrics_"
)

// 物件池：減少高併發下的 GC 壓力
var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func newContext(ctx *gin.Context) Context {
	c := contextPool.Get().(*context)
	c.ctx = ctx
	return c
}

func releaseContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil // 釋放 Gin Context
	contextPool.Put(c)
}

// 確保 context 結構實作了 Context 介面
var _ Context = (*context)(nil)

type Context interface {
	init() // 內部初始化，用於緩存 Body

	// 參數綁定系列
	ShouldBindQuery(obj interface{}) error
	ShouldBindPostForm(obj interface{}) error
	ShouldBindForm(obj interface{}) error
	ShouldBindJSON(obj interface{}) error
	ShouldBindURI(obj interface{}) error

	// 回傳控制
	Payload(payload interface{})      // 暫存成功結果
	getPayload() interface{}          // 供核心取出結果
	AbortWithError(err BusinessError) // 暫存錯誤並中斷
	abortError() BusinessError        // 供核心取出錯誤

	// 日誌與追蹤
	Logger() *zap.Logger
	setLogger(logger *zap.Logger)

	// HTTP 基礎操作
	Header() http.Header
	GetHeader(key string) string
	SetHeader(key, value string)
	Method() string
	Host() string
	Path() string
	URI() string
	RawData() []byte

	// 流程控制
	Next()
	Abort()
}

type StdContext struct {
	stdctx.Context
	*zap.Logger
}

type context struct {
	ctx *gin.Context
}

// --- 實作開始 ---

func (c *context) init() {
	// 緩存 Body 供日誌或多次讀取使用
	body, err := c.ctx.GetRawData()
	if err != nil {
		return
	}
	c.ctx.Set(_BodyName, body)
	// 重新灌回 Body，否則後續 ShouldBindJSON 會讀不到
	c.ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
}

// ShouldBind 系列封裝 Gin 的綁定功能
func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}
func (c *context) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

func (c *context) Payload(payload interface{}) {
	c.ctx.Set(_PayloadName, payload)
}

func (c *context) getPayload() interface{} {
	if payload, ok := c.ctx.Get(_PayloadName); ok {
		return payload
	}
	return nil
}

func (c *context) AbortWithError(err BusinessError) {
	if err != nil {
		httpCode := err.HTTPCode()
		if httpCode == 0 {
			httpCode = http.StatusInternalServerError
		}

		c.ctx.AbortWithStatus(httpCode)
		c.ctx.Set(_AbortErrorName, err)
	}
}

func (c *context) abortError() BusinessError {
	if err, ok := c.ctx.Get(_AbortErrorName); ok {
		return err.(BusinessError)
	}
	return nil
}

func (c *context) Logger() *zap.Logger {
	logger, ok := c.ctx.Get(_LoggerName)
	if !ok {
		return nil
	}

	return logger.(*zap.Logger)
}

func (c *context) setLogger(logger *zap.Logger) { c.ctx.Set(_LoggerName, logger) }

func (c *context) Header() http.Header {
	header := c.ctx.Request.Header

	clone := make(http.Header, len(header))
	for k, v := range header {
		value := make([]string, len(v))
		copy(value, v)

		clone[k] = value
	}
	return clone
}
func (c *context) GetHeader(key string) string { return c.ctx.GetHeader(key) }
func (c *context) SetHeader(key, value string) { c.ctx.Header(key, value) }

// Method 请求的method
func (c *context) Method() string { return c.ctx.Request.Method }

// Host 请求的host
func (c *context) Host() string { return c.ctx.Request.Host }

// Path 请求的路径(不附带querystring)
func (c *context) Path() string { return c.ctx.Request.URL.Path }

// URI unescape后的uri
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}

// RequestContext (包装 Trace + Logger) 获取请求的 context (当client关闭后，会自动canceled)
func (c *context) RequestContext() StdContext {
	return StdContext{
		//c.ctx.Request.Context(),
		stdctx.Background(),
		c.Logger(),
	}
}

// ResponseWriter 获取 ResponseWriter
func (c *context) ResponseWriter() gin.ResponseWriter {
	return c.ctx.Writer
}

func (c *context) RawData() []byte {
	if body, ok := c.ctx.Get(_BodyName); ok {
		return body.([]byte)
	}
	return nil
}

func (c *context) Next()  { c.ctx.Next() }
func (c *context) Abort() { c.ctx.Abort() }

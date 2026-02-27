package core

import (
	"net/http"
	"time"

	"github.com/brucechen520/kuji-go/internal/code"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HandlerFunc 定義
type HandlerFunc func(c Context)

type Option func(*option)

type option struct {
	enableCors bool
	// 這裡未來可以擴充 AlertNotify, RateLimit 等
}

func WithEnableCors() Option {
	return func(opt *option) {
		opt.enableCors = true
	}
}

// RouterGroup 介面封裝
type RouterGroup interface {
	Group(string, ...HandlerFunc) RouterGroup
	IRoutes
}

type IRoutes interface {
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	// 其他 Method 可依此類推...
}

type router struct {
	group *gin.RouterGroup
}

func (r *router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	group := r.group.Group(relativePath, wrapHandlers(handlers...)...)
	return &router{group: group}
}

func (r *router) GET(relativePath string, handlers ...HandlerFunc) {
	r.group.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) POST(relativePath string, handlers ...HandlerFunc) {
	r.group.POST(relativePath, wrapHandlers(handlers...)...)
}

// 核心：將我們的 HandlerFunc 轉為 Gin 的 HandlerFunc
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		h := handler // 避免閉包變數捕獲問題
		funcs[i] = func(c *gin.Context) {
			ctx := newContext(c)
			// 這裡可以加入 defer 處理 Payload 或 Error 回傳
			h(ctx)
		}
	}
	return funcs
}

type Mux interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

type mux struct {
	engine *gin.Engine
}

func (m *mux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	m.engine.ServeHTTP(w, req)
}

func (m *mux) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{
		group: m.engine.Group(relativePath, wrapHandlers(handlers...)...),
	}
}

// New 核心初始化
func New(logger *zap.Logger, options ...Option) (Mux, error) {
	gin.SetMode(gin.ReleaseMode)

	m := &mux{
		engine: gin.New(),
	}

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	// 基礎中間件：Logger, Recovery 和 統一回傳處理
	m.engine.Use(func(c *gin.Context) {
		ts := time.Now()
		ctx := newContext(c)

		// 這裡暫時手動注入 Logger，未來可進化
		ctx.setLogger(logger)

		defer func() {
			// 處理 Panic
			if err := recover(); err != nil {
				logger.Error("panic recovered", zap.Any("err", err))
				ctx.AbortWithError(Error(
					http.StatusInternalServerError,
					code.ServerError,
					code.Text(code.ServerError)),
				)
			}

			// 處理 Abort 的錯誤回傳
			if c.IsAborted() {
				if err := ctx.abortError(); err != nil {
					c.JSON(err.HTTPCode(), struct {
						Code int    `json:"code"`
						Msg  string `json:"msg"`
					}{
						Code: err.BusinessCode(),
						Msg:  err.Message(),
					})
				}
			} else {
				// 正確回傳 Payload
				if payload := ctx.getPayload(); payload != nil {
					c.JSON(http.StatusOK, payload)
				}
			}

			logger.Info("request access",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Duration("cost", time.Since(ts)),
			)
		}()

		c.Next()
	})

	return m, nil
}

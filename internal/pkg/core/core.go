package core

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/brucechen520/kuji-go/pkg/color"
	"github.com/brucechen520/kuji-go/pkg/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const _UI = `
██╗  ██╗██╗   ██╗     ██╗██╗      ██████╗  ██████╗
██║ ██╔╝██║   ██║     ██║██║     ██╔════╝ ██╔═══██╗
█████╔╝ ██║   ██║     ██║██║     ██║  ███╗██║   ██║
██╔═██╗ ██║   ██║██   ██║██║     ██║   ██║██║   ██║
██║  ██╗╚██████╔╝╚██████╔╝██║     ╚██████╔╝╚██████╔╝
╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚═╝      ╚═════╝  ╚═════╝
`

type Option func(*option)

type option struct {
	disablePProf      bool
	disableSwagger    bool
	disablePrometheus bool
	enableCors        bool
	enableRate        bool
	enableOpenBrowser string
	// alertNotify       proposal.NotifyHandler
	// recordHandler     proposal.RecordHandler
}

// WithDisablePProf 禁用 pprof
func WithDisablePProf() Option {
	return func(opt *option) {
		opt.disablePProf = true
	}
}

// WithDisableSwagger 禁用 swagger
func WithDisableSwagger() Option {
	return func(opt *option) {
		opt.disableSwagger = true
	}
}

// WithDisablePrometheus 禁用prometheus
func WithDisablePrometheus() Option {
	return func(opt *option) {
		opt.disablePrometheus = true
	}
}

// WithAlertNotify 设置告警通知
// TODO: 目前先保留接口，后续实现告警通知功能
// func WithAlertNotify(notifyHandler proposal.NotifyHandler) Option {
// 	return func(opt *option) {
// 		opt.alertNotify = notifyHandler
// 	}
// }

// WithRecordMetrics 设置记录接口指标
// TODO: 目前先保留接口，后续实现记录接口指标功能
// func WithRecordMetrics(recordHandler proposal.RecordHandler) Option {
// 	return func(opt *option) {
// 		opt.recordHandler = recordHandler
// 	}
// }

// WithEnableOpenBrowser 启动后在浏览器中打开 uri
func WithEnableOpenBrowser(uri string) Option {
	return func(opt *option) {
		opt.enableOpenBrowser = uri
	}
}

// WithEnableCors 设置支持跨域
func WithEnableCors() Option {
	return func(opt *option) {
		opt.enableCors = true
	}
}

// WithEnableRate 设置支持限流
func WithEnableRate() Option {
	return func(opt *option) {
		opt.enableRate = true
	}
}

// RouterGroup 定義路由群組介面，支援鏈式調用
type RouterGroup interface {
	Group(string, ...HandlerFunc) RouterGroup
	IRoutes
}

var _ IRoutes = (*router)(nil)

// IRoutes 包装gin的IRoutes
type IRoutes interface {
	Any(string, ...HandlerFunc)
	GET(string, ...HandlerFunc)
	POST(string, ...HandlerFunc)
	DELETE(string, ...HandlerFunc)
	PATCH(string, ...HandlerFunc)
	PUT(string, ...HandlerFunc)
	OPTIONS(string, ...HandlerFunc)
	HEAD(string, ...HandlerFunc)
}

// router 實作 RouterGroup 介面
type router struct {
	group *gin.RouterGroup
}

func (r *router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	group := r.group.Group(relativePath, wrapHandlers(handlers...)...)
	return &router{group: group}
}

func (r *router) Any(relativePath string, handlers ...HandlerFunc) {
	r.group.Any(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) GET(relativePath string, handlers ...HandlerFunc) {
	r.group.GET(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) POST(relativePath string, handlers ...HandlerFunc) {
	r.group.POST(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.group.DELETE(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PATCH(relativePath string, handlers ...HandlerFunc) {
	r.group.PATCH(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) PUT(relativePath string, handlers ...HandlerFunc) {
	r.group.PUT(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) OPTIONS(relativePath string, handlers ...HandlerFunc) {
	r.group.OPTIONS(relativePath, wrapHandlers(handlers...)...)
}

func (r *router) HEAD(relativePath string, handlers ...HandlerFunc) {
	r.group.HEAD(relativePath, wrapHandlers(handlers...)...)
}

// wrapHandlers 將我們自定義的 HandlerFunc 轉為 gin.HandlerFunc
func wrapHandlers(handlers ...HandlerFunc) []gin.HandlerFunc {
	funcs := make([]gin.HandlerFunc, len(handlers))
	for i, handler := range handlers {
		handler := handler // 避免閉包變數捕獲問題
		funcs[i] = func(c *gin.Context) {
			ctx := newContext(c) // 從池子裡拿 Context
			defer releaseContext(ctx)
			handler(ctx)
		}
	}
	return funcs
}

// Mux http mux
type Mux interface {
	http.Handler
	Group(relativePath string, handlers ...HandlerFunc) RouterGroup
}

// mux 實作 Mux 介面
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

func New(logger *zap.Logger, options ...Option) (Mux, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	gin.SetMode(gin.ReleaseMode)
	mux := &mux{
		engine: gin.New(),
	}

	fmt.Println(color.Blue(_UI))

	// TODO: 載入 靜態資源與模板：

	// TODO: withoutTracePaths 这些请求，默认不记录日志

	opt := new(option)
	for _, f := range options {
		f(opt)
	}

	// TODO: 根據 opt 的設定，註冊對應的 Middleware 和路由
	// TODO: if !opt.disablePProf {
	// TODO: if !opt.disableSwagger {
	// TODO: if !opt.disablePrometheus {
	// TODO: if opt.enableCors {
	// TODO: if opt.enableOpenBrowser != "" {

	// 1. recover两次，防止处理时发生panic，尤其是在OnPanicNotify中。
	mux.engine.Use(func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("got panic", zap.String("panic", fmt.Sprintf("%+v", err)), zap.String("stack", string(debug.Stack())))
			}
		}()

		ctx.Next()
	})

	mux.engine.Use(func(ctx *gin.Context) {
		// 如果已經 404，直接跳過不處理
		if ctx.Writer.Status() == http.StatusNotFound {
			return
		}

		ts := time.Now() // 開始計時

		// 1. 從 Pool 取得自定義 Context 並初始化
		c := newContext(ctx)
		defer releaseContext(c)

		c.init()
		c.setLogger(logger)

		// 2. 設定 defer 處理請求結束後的「善後」
		defer func() {
			// --- 第二次 Recover：處理業務 Panic ---
			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				logger.Error("business panic",
					zap.Any("error", err),
					zap.String("stack", stackInfo),
				)

				// TODO: notifyHandler - 發送告警通知 (Slack/Email)

				// Panic 時確保回傳 500
				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code":    500,
					"message": "Internal Server Error",
				})
			}

			// --- 處理錯誤回傳 ---
			if ctx.IsAborted() {
				if err := c.abortError(); err != nil {
					// TODO: 這裡可以根據自定義錯誤決定是否發送告警
					ctx.JSON(err.HTTPCode(), map[string]interface{}{
						"code":    err.BusinessCode(),
						"message": err.Message(),
					})
					return
				}
			}

			// --- 處理成功回傳 ---
			// 注意：這裡先保留你原本的邏輯，統一回傳 200
			response := c.getPayload()
			if response != nil && !ctx.IsAborted() {
				ctx.JSON(http.StatusOK, response)
			}

			// --- 最終日誌紀錄與耗時統計 ---
			costSeconds := time.Since(ts).Seconds()

			// TODO: 記錄指標 (Metrics) 到 Prometheus

			logger.Info("request-log",
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.Request.URL.Path),
				zap.Int("status", ctx.Writer.Status()),
				zap.Float64("cost_seconds", costSeconds),
				// TODO: 紀錄 TraceID
			)
		}()

		ctx.Next() // 執行後續 Handler
	})

	// TODO: Rate Limiter 限流機制
	// 在進入耗時的邏輯前，先擋掉超過頻率的請求

	// TODO: TraceID 鏈路追蹤注入
	// 在這裡生成 TraceID，後面的所有 Log 才能共享同一個 ID

	// 第二層：Context 注入、耗時統計、統一回覆
	mux.engine.Use(func(ctx *gin.Context) {
		ts := time.Now()

		c := newContext(ctx)
		defer releaseContext(c)

		c.init()
		c.setLogger(logger)

		// 第二層 Recover：處理業務邏輯 Panic 並格式化輸出
		defer func() {
			if err := recover(); err != nil {
				stackInfo := string(debug.Stack())
				logger.Error("business panic",
					zap.Any("error", err),
					zap.String("stack", stackInfo),
				)

				ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
					"code":    500,
					"message": "Internal Server Error (Panic)",
				})
			}

			// 統一回覆機制
			response := c.getPayload()
			if response != nil && !ctx.IsAborted() {
				ctx.JSON(http.StatusOK, response)
			}

			// 耗時紀錄 (如果有 TraceID TODO，這裡記得補上 TraceID)
			costSeconds := time.Since(ts).Seconds()
			logger.Info("request-log",
				zap.String("method", ctx.Request.Method),
				zap.String("path", ctx.Request.URL.Path),
				zap.Int("status", ctx.Writer.Status()),
				zap.Float64("cost_seconds", costSeconds),
			)
		}()

		ctx.Next()
	})

	// 系統預設路由, TODO: 上面要＋ health, metrics, swagger,pprof 等 routes 不記錄 logger
	system := mux.Group("/system")
	{
		// 健康检查
		system.GET("/health", func(ctx Context) {
			resp := &struct {
				Timestamp   time.Time `json:"timestamp"`
				Environment string    `json:"environment"`
				Host        string    `json:"host"`
				Status      string    `json:"status"`
			}{
				Timestamp:   time.Now(),
				Environment: "development", // TODO: 這裡可以根據環境變數動態設定
				Host:        ctx.Host(),
				Status:      "ok",
			}
			ctx.Payload(resp)
		})
	}

	return mux, nil
}

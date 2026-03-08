package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	cors "github.com/rs/cors/wrapper/gin"
	"go.uber.org/zap"
)

// Engine 封裝了 gin.Engine，提供初始化後的引擎實例
type Engine struct {
	*gin.Engine
}

type Option func(*option)

type option struct {
	disablePProf      bool
	disableSwagger    bool
	disablePrometheus bool
	enableCors        bool
	enableRate        bool
	// enableOpenBrowser string
	// alertNotify       proposal.NotifyHandler
	// recordHandler     proposal.RecordHandler
}

// NewEngine 透過注入 Logger，組裝一個完整的 Gin 引擎
// 這裡符合 Wire 的 Provider 定義
func NewEngine(l *zap.Logger, opts ...Option) *Engine {
	// 這裡設定 gin 為生產模式，減少開發模式的額外輸出
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 註冊我們自定義的 Middleware
	// 這裡呼叫了 middleware.go 中的邏輯
	r.Use(LoggerMiddleware(l))

	// 也可以在這裡註冊 Recover 等其他共用 Middleware
	r.Use(gin.Recovery())

	e := &Engine{r}

	opt := new(option)
	// 執行所有外掛選項 (Cors, RateLimit 等)
	for _, f := range opts {
		f(opt)
	}

	if !opt.disableSwagger {
		e.GET("/swagger/*any", nil) // register swagger
	}

	if !opt.disablePrometheus {
		e.GET("/metrics", gin.WrapH(promhttp.Handler())) // register prometheus
	}

	if opt.enableCors {
		e.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:     []string{"*"},
			AllowCredentials:   true,
			OptionsPassthrough: true,
		}))
	}

	// TODO: other options

	return &Engine{r}
}

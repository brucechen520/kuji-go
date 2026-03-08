package core

import (
	"net/http"
	"net/url"
	"time"

	"github.com/brucechen520/kuji-go/pkg/trace"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// withoutTracePaths 这些请求，默认不记录日志
var withoutTracePaths = map[string]bool{
	"/metrics": true,

	"/debug/pprof/":             true,
	"/debug/pprof/cmdline":      true,
	"/debug/pprof/profile":      true,
	"/debug/pprof/symbol":       true,
	"/debug/pprof/trace":        true,
	"/debug/pprof/allocs":       true,
	"/debug/pprof/block":        true,
	"/debug/pprof/goroutine":    true,
	"/debug/pprof/heap":         true,
	"/debug/pprof/mutex":        true,
	"/debug/pprof/threadcreate": true,

	"/favicon.ico": true,

	"/system/health": true,
}

func LoggerMiddleware(l *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ts := time.Now()

		customCtx := NewContext(c)
		defer ReleaseContext(customCtx)

		// 初始化 context, 包括讀取 raw data
		customCtx.(*customContext).init()

		if !withoutTracePaths[c.Request.URL.Path] {
			if traceId := c.GetHeader(trace.Header); traceId != "" {
				customCtx.setTrace(trace.New(traceId))
			} else {
				customCtx.setTrace(trace.New(""))
			}
		}

		// 2. 建立 Trace 與 Logger
		reqLogger := l.With(zap.String("trace_id", customCtx.Trace().ID()))

		// 3. 注入
		customCtx.setLogger(reqLogger)

		c.Set(trace.Header, customCtx.Trace().ID())
		defer func() {
			decodedURL, _ := url.QueryUnescape(c.Request.URL.RequestURI())

			// ctx.Request.Header，精简 Header 参数
			traceHeader := map[string]string{
				"Content-Type": c.GetHeader("Content-Type"),
				// configs.HeaderLoginToken:    c.GetHeader(configs.HeaderLoginToken),
				// configs.HeaderSignToken:     c.GetHeader(configs.HeaderSignToken),
				// configs.HeaderSignTokenDate: c.GetHeader(configs.HeaderSignTokenDate),
			}

			// region 记录日志
			var t *trace.Trace
			if x := customCtx.Trace(); x != nil {
				// customCtx.Trace() 是回傳 trace.T 這個 interface，所以要轉型
				t = x.(*trace.Trace)
			} else {
				return
			}

			// 請求結束後，一次性輸出 Trace 總結
			t.WithRequest(&trace.Request{
				TTL:        "un-limit",
				Method:     c.Request.Method,
				DecodedURL: decodedURL,
				Header:     traceHeader,
				Body:       string(customCtx.RawData()),
			})

			// var responseBody interface{}

			// if response != nil {
			// 	responseBody = response
			// }

			// graphResponse = context.getGraphPayload()
			// if graphResponse != nil {
			// 	responseBody = graphResponse
			// }

			// t.WithResponse(&trace.Response{
			// 	Header:          ctx.Writer.Header(),
			// 	HttpCode:        ctx.Writer.Status(),
			// 	HttpCodeMsg:     http.StatusText(ctx.Writer.Status()),
			// 	BusinessCode:    businessCode,
			// 	BusinessCodeMsg: businessCodeMsg,
			// 	Body:            responseBody,
			// 	CostSeconds:     time.Since(ts).Seconds(),
			// })

			t.Success = !c.IsAborted() && (c.Writer.Status() == http.StatusOK)
			t.CostSeconds = time.Since(ts).Seconds()
			customCtx.GetLogger().Info("request finished",
				zap.Any("method", c.Request.Method),
				zap.Any("path", decodedURL),
				zap.Any("http_code", c.Writer.Status()),
				// zap.Any("business_code", businessCode),
				zap.Any("success", t.Success),
				zap.Any("cost_seconds", t.CostSeconds),
				zap.Any("trace_id", t.Identifier),
				zap.Any("trace_info", t),
				// zap.Error(abortErr),
			)
			// endregion
		}()

		c.Next()
	}
}

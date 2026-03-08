package logger

import (
	"os"

	"github.com/brucechen520/kuji-go/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// option 定義 Logger 行為與組件開關
type option struct {
	level        zapcore.Level
	redisEnabled bool
	dbEnabled    bool
	skipPaths    map[string]struct{}
}

// Option 定義功能選項函式
type Option func(*option)

// NewLogger 初始化 Logger 工廠
func NewLogger(opts ...Option) *zap.Logger {
	cfg := &option{
		level:     zap.InfoLevel,
		skipPaths: make(map[string]struct{}),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	// 這裡示範初始化生產環境用的 JSON Logger
	encoderConfig := zap.NewProductionEncoderConfig()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		cfg.level,
	)

	return zap.New(core)
}

// Option 建構函式
func WithLevel(level zapcore.Level) Option {
	return func(c *option) { c.level = level }
}

func WithRedisLog(enabled bool) Option {
	return func(c *option) { c.redisEnabled = enabled }
}

func WithSkipPaths(paths []string) Option {
	return func(c *option) {
		for _, p := range paths {
			c.skipPaths[p] = struct{}{}
		}
	}
}

// NewRequestLogger 專門給 Gin Middleware 使用
func NewRequestLogger(cfg *config.Config) *zap.Logger {
	// 這裡可以針對 Gin 請求日誌設定特定的 Level 或輸出
	return NewLogger(WithLevel(zap.InfoLevel))
}

// NewWorkerLogger 預留給未來使用
func NewWorkerLogger(cfg *config.Config) *zap.Logger {
	return NewLogger(WithLevel(zap.DebugLevel))
}

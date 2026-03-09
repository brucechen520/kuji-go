package trace

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gormTraceLogger 實作了 gorm.io/gorm/logger.Interface。
// 透過工廠函數 NewGormTraceLogger 取得，外部只依賴 logger.Interface，不直接接觸此 struct。
type gormTraceLogger struct {
	zapLogger                 *zap.Logger
	logLevel                  logger.LogLevel
	slowThreshold             time.Duration
	ignoreRecordNotFoundError bool
}

// NewGormTraceLogger 建立一個負責追蹤 SQL 的 GORM Logger
func NewGormTraceLogger(zapLogger *zap.Logger, config logger.Config) logger.Interface {
	return &gormTraceLogger{
		zapLogger:                 zapLogger,
		logLevel:                  config.LogLevel,
		slowThreshold:             config.SlowThreshold,
		ignoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
	}
}

// LogMode 實作 gorm logger.Interface
func (l *gormTraceLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.logLevel = level
	return &newLogger
}

// Info 實作 gorm logger.Interface
func (l *gormTraceLogger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.logLevel >= logger.Info {
		l.zapLogger.Sugar().Infof(s, args...)
	}
}

// Warn 實作 gorm logger.Interface
func (l *gormTraceLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.logLevel >= logger.Warn {
		l.zapLogger.Sugar().Warnf(s, args...)
	}
}

// Error 實作 gorm logger.Interface
func (l *gormTraceLogger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.logLevel >= logger.Error {
		l.zapLogger.Sugar().Errorf(s, args...)
	}
}

// Trace 負責攔截與記錄執行完畢的 SQL 軌跡 (最重要的方法)
func (l *gormTraceLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.logLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 1. 產生標準化日誌 (Zap)
	switch {
	case err != nil && l.logLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.ignoreRecordNotFoundError):
		l.zapLogger.Error("SQL execution error",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.logLevel >= logger.Warn:
		l.zapLogger.Warn("Slow SQL query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.logLevel == logger.Info:
		l.zapLogger.Info("SQL query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}

	// 2. 將日誌掛載至 Context Trace 中 (重點機制)
	// 透過我們設計的 ExtractTrace 由標準 context 取回 Trace 物件
	if traceObj := ExtractTrace(ctx); traceObj != nil {
		traceObj.AppendSQL(&SQL{
			Timestamp:   begin.Format("2006-01-02 15:04:05.000"),
			Stack:       "",
			SQL:         sql,
			Rows:        rows,
			CostSeconds: elapsed.Seconds(),
		})
	}
}

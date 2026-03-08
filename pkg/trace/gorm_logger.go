package trace

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormTraceLogger 實作了 gorm.io/gorm/logger.Interface
type GormTraceLogger struct {
	ZapLogger                 *zap.Logger
	LogLevel                  logger.LogLevel
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
}

// NewGormTraceLogger 建立一個負責追蹤 SQL 的 GORM Logger
func NewGormTraceLogger(zapLogger *zap.Logger, config logger.Config) logger.Interface {
	return &GormTraceLogger{
		ZapLogger:                 zapLogger,
		LogLevel:                  config.LogLevel,
		SlowThreshold:             config.SlowThreshold,
		IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
	}
}

// LogMode 實作 gorm logger.Interface
func (l *GormTraceLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 實作 gorm logger.Interface
func (l *GormTraceLogger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Sugar().Infof(s, args...)
	}
}

// Warn 實作 gorm logger.Interface
func (l *GormTraceLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Sugar().Warnf(s, args...)
	}
}

// Error 實作 gorm logger.Interface
func (l *GormTraceLogger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Sugar().Errorf(s, args...)
	}
}

// Trace 負責攔截與記錄執行完畢的 SQL 軌跡 (最重要的方法)
func (l *GormTraceLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 1. 產生標準化日誌 (Zap)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.ZapLogger.Error("SQL execution error",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.ZapLogger.Warn("Slow SQL query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.LogLevel == logger.Info:
		l.ZapLogger.Info("SQL query",
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}

	// 2. 將日誌掛載至 Context Trace 中 (重點機制)
	// 透過我們設計的 ExtractTrace 由標準 context 取回 Trace 物件
	if traceObj := ExtractTrace(ctx); traceObj != nil {
		traceObj.AppendSQL(&SQL{
			Timestamp:   begin.Format("2006-01-02 15:04:05.000"), // 加上毫秒比較好追蹤
			Stack:       "",                                      // Gorm 複雜的 caller 抓取這邊為了效能暫免
			SQL:         sql,
			Rows:        rows,
			CostSeconds: elapsed.Seconds(),
		})
	}
}

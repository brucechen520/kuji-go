package trace

import (
	"context"
	"net"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisTraceHook 實作了 go-redis/v9 的 redis.Hook 介面
type RedisTraceHook struct {
	ZapLogger *zap.Logger
}

func NewRedisTraceHook(zapLogger *zap.Logger) redis.Hook {
	return &RedisTraceHook{
		ZapLogger: zapLogger,
	}
}

// DialHook 攔截連線建立 (我們這裡不修改它，直接呼叫下一個)
func (h *RedisTraceHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook 攔截單一指令執行
func (h *RedisTraceHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmd) // 真正執行指令

		elapsed := time.Since(start)

		// 格式化指令
		cmdName := cmd.Name()
		cmdString := cmd.String()

		if err != nil && err != redis.Nil {
			h.ZapLogger.Error("Redis execute error",
				zap.Error(err),
				zap.String("cmd", cmdString),
				zap.Duration("elapsed", elapsed),
			)
		} else {
			h.ZapLogger.Debug("Redis execute",
				zap.String("cmd", cmdString),
				zap.Duration("elapsed", elapsed),
			)
		}

		// 掛載追蹤資料到 Context Trace
		if traceObj := ExtractTrace(ctx); traceObj != nil {
			// 將 key 和 value 從字串拆解出來
			// cmdString 格式通常為 "get key:name" 或 "set key value"
			parts := strings.SplitN(cmdString, " ", 3)
			key := ""
			val := ""
			if len(parts) >= 2 {
				key = parts[1]
			}
			if len(parts) == 3 {
				val = parts[2]
			}

			traceObj.AppendRedis(&Redis{
				Timestamp:   start.Format("2006-01-02 15:04:05.000"),
				Handle:      strings.ToUpper(cmdName),
				Key:         key,
				Value:       val,
				CostSeconds: elapsed.Seconds(),
			})
		}

		return err
	}
}

// ProcessPipelineHook 攔截 Pipeline 多重指令執行
func (h *RedisTraceHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		start := time.Now()

		err := next(ctx, cmds)

		elapsed := time.Since(start)

		h.ZapLogger.Debug("Redis pipeline execute",
			zap.Int("cmd_count", len(cmds)),
			zap.Duration("elapsed", elapsed),
		)

		if traceObj := ExtractTrace(ctx); traceObj != nil {
			traceObj.AppendRedis(&Redis{
				Timestamp:   start.Format("2006-01-02 15:04:05.000"),
				Handle:      "PIPELINE",
				Key:         "multiple",
				CostSeconds: elapsed.Seconds(),
			})
		}

		return err
	}
}

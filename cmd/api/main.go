package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"go.uber.org/zap"

	"github.com/brucechen520/kuji-go/internal/config"
)

func main() {
	// 建立一個臨時的基礎 logger 供初始化階段使用
	baseLogger, _ := zap.NewProduction()
	defer baseLogger.Sync() // flushes buffer, if any

	// 1. 載入設定 (從環境變數或 CD 注入的設定檔)
	cfg, err := config.Load()
	if err != nil {
		baseLogger.Fatal("無法載入設定檔", zap.Error(err))
	}

	// 2. 透過 Wire 產生的函數進行依賴注入
	// router: Gin Engine, cleanup: 關閉 DB/Redis 的函數
	router, cleanup, err := InitializeApp(cfg)
	if err != nil {
		baseLogger.Fatal("依賴注入初始化失敗", zap.Error(err))
	}
	defer func() {
		cleanup()
		baseLogger.Info("Close redis and database connection")
	}() // 程式結束時確保連線資源釋放

	// 3. 設定 HTTP Server (為了支援優雅關閉，不直接用 router.Run)
	srv := &http.Server{
		Addr:    ":" + cfg.App.Port, // 假設你在 config 有定義 Port
		Handler: router,
	}

	// 4. 在 Goroutine 中啟動 Server，避免阻塞主執行緒
	go func() {
		baseLogger.Info("服務啟動中...", zap.String("port", cfg.App.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			baseLogger.Fatal("監聽失敗", zap.Error(err))
		}
	}()

	// 5. 等待中斷訊號 (優雅關閉 Graceful Shutdown)
	quit := make(chan os.Signal, 1)
	// 監聽 Ctrl+C (SIGINT) 或 系統刪除訊號 (SIGTERM)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	baseLogger.Info("正在關閉伺服器...")

	// 設定 5 秒超時，給予正在處理的請求緩衝時間
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		baseLogger.Info("緩衝時間結束關閉伺服器中")
	}()

	if err := srv.Shutdown(ctx); err != nil {
		baseLogger.Fatal("伺服器強制關閉", zap.Error(err))
	}

	baseLogger.Info("伺服器已安全退出")
}

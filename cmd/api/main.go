package main

import (
	"context"
	"net/http"
	"time"

	"github.com/brucechen520/kuji-go/internal/router"
	"github.com/brucechen520/kuji-go/pkg/shutdown"
	"github.com/joho/godotenv"

	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load() // 自動讀取 .env 並寫入 os.Environ

	// 1. 初始化 Zap Logger (取代原本複雜的 logger.NewJSONLogger)
	// 在開發環境建議用 NewDevelopment，生產環境用 NewProduction
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // 確保程式結束前日誌有刷入磁碟

	// 2. 初始化 HTTP 服務 (組裝所有 Repo, Service, Handler)
	// 我們把原本的兩個 logger 簡化為一個傳進去
	s, err := router.NewHTTPServer(logger)
	if err != nil {
		logger.Fatal("HTTP Server 初始化失敗", zap.Error(err))
	}

	// 3. 設定 HTTP Server
	server := &http.Server{
		Addr:    "127.0.0.1:8080", // 確保 configs 裡有定義 Port
		Handler: s.Mux,
	}

	// 4. 啟動服務 (非阻塞)
	go func() {
		logger.Info("服務啟動中...", zap.String("port", "8080"))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP Server 啟動異常", zap.Error(err))
		}
	}()

	// 5. 優雅關閉 Hook
	shutdown.NewHook().Close(
		// 關閉 HTTP Server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				logger.Error("HTTP Server 關閉失敗", zap.Error(err))
			}
		},

		// 關閉 DB 連線
		func() {
			if s.Db != nil {
				// 這裡調用你內部的關閉邏輯
				if err := s.Db.DbWClose(); err != nil {
					logger.Error("資料庫(W)關閉失敗", zap.Error(err))
				}
			}
		},

		// 關閉快取 (Redis)
		func() {
			if s.Rdb != nil {
				if err := s.Rdb.RDBWClose(); err != nil {
					logger.Error("快取服務關閉失敗", zap.Error(err))
				}
			}
		},
	)
}

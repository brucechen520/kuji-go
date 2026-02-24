package main // main 套件是 Go 程式的執行入口

import (
	"context"                 // 引入 context
	"kuji-go/internal/app"    // 引入 app 套件，負責應用程式組裝
	"kuji-go/internal/router" // 引入 router 套件
	"log"                     // 引入標準日誌套件
	"net/http"                // 引入 http 套件
	"os"                      // 引入 os 套件
	"os/signal"               // 引入 signal 套件
	"syscall"                 // 引入 syscall
	"time"                    // 引入 time 套件
)

// main 函式是程式執行的起點
func main() {
	// 1. 初始化應用程式容器
	// 呼叫 app.NewContainer() 完成所有依賴的組裝 (DB -> Repo -> Service -> Handler)
	// 接收 cleanup 函式，準備在程式結束時執行
	container, cleanup := app.NewContainer()

	// 確保在 main 函式結束前 (無論是正常結束還是 panic) 關閉資料庫與 Redis
	defer cleanup()

	// 2. 設定路由
	// 從容器中取出組裝好的 Handler 傳給 Router
	r := router.SetupRouter(container.Handler)

	// 3. 啟動伺服器
	// 使用 http.Server 取代 r.Run()，以便控制關閉流程
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 在 Goroutine 中啟動伺服器，避免阻塞主線程
	go func() {
		log.Println("一番賞系統成功啟動於 :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("伺服器啟動失敗: %s\n", err)
		}
	}()

	// 4. 優雅關閉 (Graceful Shutdown)
	// 等待中斷信號 (如 Ctrl+C 或 Docker stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞直到收到信號
	log.Println("正在關閉伺服器...")

	// 設定 5 秒的超時時間，讓正在處理的請求有時間完成
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 呼叫 Shutdown 停止接收新請求並等待舊請求完成
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("伺服器強制關閉:", err)
	}

	log.Println("伺服器已優雅關閉")
	// 這裡會自動執行 defer cleanup()，關閉 DB 和 Redis
}

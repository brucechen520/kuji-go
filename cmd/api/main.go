package main // main 套件是 Go 程式的執行入口

import (
	"kuji-go/internal/app"    // 引入 app 套件，負責應用程式組裝
	"kuji-go/internal/router" // 引入 router 套件
	"log"                     // 引入標準日誌套件
)

// main 函式是程式執行的起點
func main() {
	// 1. 初始化應用程式容器
	// 呼叫 app.NewContainer() 完成所有依賴的組裝 (DB -> Repo -> Service -> Handler)
	container := app.NewContainer()

	// 2. 設定路由
	// 從容器中取出組裝好的 Handler 傳給 Router
	r := router.SetupRouter(container.Handler)

	// 3. 啟動伺服器
	log.Println("一番賞系統成功啟動於 :8080") // 印出啟動訊息
	if err := r.Run(":8080"); err != nil {
		log.Fatal("伺服器啟動失敗: ", err) // r.Run 預設監聽 8080 port，若失敗則終止
	}
}

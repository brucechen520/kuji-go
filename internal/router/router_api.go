package router

import (
	prize_handler "github.com/brucechen520/kuji-go/internal/handlers/prize"
	"github.com/brucechen520/kuji-go/internal/repository"
	prize_service "github.com/brucechen520/kuji-go/internal/service/client/prize"
)

func setApiRouter(r *resource) {
	// --- 這裡就是你的 DI 組裝區 ---

	// 1. Repository 層 (注入 DB 和 Redis)
	prizeRepo := repository.NewPrizeRepo(r.db, r.rdb)

	// 2. Service 層 (注入 Repository，如果需要事務管理也可以注入 transactionManager)
	prizeSrv := prize_service.New(prizeRepo, r.transactionManager)

	// 3. Handler 層 (注入 Service)
	prizeHandler := prize_handler.New(prizeSrv, r.logger)

	// --- 路由定義區 ---
	// 使用 r.mux 進行路由群組與路徑綁定
	v1 := r.mux.Group("/api/v1")
	{
		client := v1.Group("/client")

		// 取得多個獎項 (複數)
		client.GET("/prizes", prizeHandler.ListPrizes)
	}
}

//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/brucechen520/kuji-go/internal/config"
	clientH "github.com/brucechen520/kuji-go/internal/handler/client" // 確保引用了 http 層
	"github.com/brucechen520/kuji-go/internal/pkg"
	"github.com/brucechen520/kuji-go/internal/repository"
	"github.com/brucechen520/kuji-go/internal/route"
	clientSrv "github.com/brucechen520/kuji-go/internal/service/client"
)

func InitializeApp(cfg *config.Config) (*gin.Engine, func(), error) {
	wire.Build(
		// 1. 設定檔拆解：告訴 Wire 如何從大 Config 拿到 AuthConfig
		// 注意：這裡的 "Auth" 必須對應到你 config.Config 結構體裡的成員名稱
		wire.FieldsOf(new(*config.Config), "Auth"),

		// 2. 注入 Repository 層 (包含 InitDB, InitRedis, NewUserRepo, NewTokenStore)
		repository.ProviderSet,

		// 3. 注入 Service 層 (建議在 service 套件也定義一個 ProviderSet)
		clientSrv.ProviderSet,

		// 4. 注入傳輸層 (Handler)
		clientH.ProviderSet,

		// 5. 注入 pkg
		pkg.ProviderSet,

		// 6. 注入 Route
		route.ProviderSet,

		// 5. 最終組裝：回傳 Gin Engine
		// Wire 會自動根據 http.NewRouter 的回傳值來滿足 InitializeApp 的第一個回傳值
	)
	return nil, nil, nil
}

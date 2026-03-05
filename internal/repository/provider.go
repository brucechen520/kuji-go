package repository

import (
	"github.com/google/wire"

	"github.com/brucechen520/kuji-go/internal/repository/postgre"
	REDIS_REPO "github.com/brucechen520/kuji-go/internal/repository/redis"
)

// ProviderSet 是一張「零件清單」，告訴 Wire 這層樓有哪些工具可以用
var ProviderSet = wire.NewSet(
	postgre.NewUserRepo,      // 這裡是你具體的 MySQL/Postgres 實作
	REDIS_REPO.NewTokenStore, // 這裡是你具體的 Redis 實作
)

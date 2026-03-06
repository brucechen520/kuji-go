package repository

import (
	"github.com/google/wire"

	"github.com/brucechen520/kuji-go/internal/repository/postgre/client"
	"github.com/brucechen520/kuji-go/internal/repository/redis"
)

// ProviderSet 是一張「零件清單」，告訴 Wire 這層樓有哪些工具可以用
var ProviderSet = wire.NewSet(
	client.NewUserRepo,
	client.NewSeriesRepo,
	redis.NewTokenStore,
	redis.NewKujiStore,
)

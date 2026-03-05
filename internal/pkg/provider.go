package pkg

import (
	"github.com/google/wire"
)

// ProviderSet 是一張「零件清單」，告訴 Wire 這層樓有哪些工具可以用
var ProviderSet = wire.NewSet(
	InitDB,
	InitRedis,
)

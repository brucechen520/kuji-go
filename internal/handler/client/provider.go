package client // 請確認這與你 wire.go import 的名稱一致

import "github.com/google/wire"

// 把這層樓所有的 Handler 都列在這裡，也包含 Router
var ProviderSet = wire.NewSet(
	NewAuthHandler,
)

package client // 請確認這與你 wire.go import 的名稱一致

import "github.com/google/wire"

// 把這層樓所有的 Service 都列在這裡
var ProviderSet = wire.NewSet(
	NewAuthService,
	NewSeriesService,
	// 如果有其他 Service 也寫在這裡
)

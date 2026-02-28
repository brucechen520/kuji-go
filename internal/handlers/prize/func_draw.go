// handlers/prize/func_draw.go
package prize

import "github.com/brucechen520/kuji-go/internal/pkg/core"

func (h *Handler) Draw() core.HandlerFunc {
	return func(ctx core.Context) {
		// 1. 驗證參數 (Request Binding)
		// 2. 呼叫 Service
		// 3. 回傳 Payload
	}
}

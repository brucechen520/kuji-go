package handlers // 屬於 handlers 套件

// Handler 是總入口，負責聚合所有子 Handler
// 這種設計模式可以讓 main.go 只需要初始化這個大 Handler，而不需要管理每個小的 Handler
type Handler struct {
	Prize *PrizeHandler // 包含 PrizeHandler 子模組
	// 未來可以加入 User *UserHandler, Box *BoxHandler 等其他模組
}

// NewHandler 初始化總 Handler
// 改為接收已經初始化好的 PrizeHandler (由 app.go 組裝)
func NewHandler(prizeH *PrizeHandler) *Handler {
	return &Handler{
		Prize: prizeH, // 設定 PrizeHandler
	}
}

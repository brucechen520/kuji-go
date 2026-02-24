package handlers // 屬於 handlers 套件

// Option is a functional option for configuring the Handler.
type Option func(*Handler)

// Handler 是總入口，負責聚合所有子 Handler
// 這種設計模式可以讓 main.go 只需要初始化這個大 Handler，而不需要管理每個小的 Handler
type Handler struct {
	Prize *PrizeHandler // 包含 PrizeHandler 子模組
	// User  *UserHandler // 未來可以加入 UserHandler
}

// NewHandler 初始化總 Handler
// 使用 Functional Options Pattern，使其符合 OCP (Open/Closed Principle)
func NewHandler(opts ...Option) *Handler {
	h := &Handler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// WithPrizeHandler returns an Option to set the PrizeHandler.
func WithPrizeHandler(ph *PrizeHandler) Option {
	return func(h *Handler) {
		h.Prize = ph
	}
}

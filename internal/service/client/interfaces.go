package client

import (
	"github.com/brucechen520/kuji-go/internal/dto"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
)

// SeriesService defines the business operations for the series domain.
// Handler 依賴此介面而非具體的 *SeriesService struct，
// 使得 Handler 層可以獨立單元測試（注入 MockSeriesService），
// 也允許未來替換實作而不影響上層。
type SeriesService interface {
	GetSeriesById(ctx core.Context, id uint) (*dto.SeriesDetailDTO, error)
}

// AuthService defines the business operations for the auth domain.
type AuthService interface {
	Login(ctx core.Context, email string) (string, error)
}

package client

import (
	"strconv"
	"time"

	"github.com/brucechen520/kuji-go/internal/config"
	"github.com/brucechen520/kuji-go/internal/pkg/core"
	"github.com/brucechen520/kuji-go/internal/repository/postgre/client"
	"github.com/brucechen520/kuji-go/internal/repository/redis"
)

// compile-time check
var _ AuthService = (*authService)(nil)

type authService struct {
	userRepo   client.UserRepository
	tokenStore redis.TokenStore
}

func NewAuthService(ur client.UserRepository, ts redis.TokenStore, cfg *config.AuthConfig) AuthService {
	return &authService{userRepo: ur, tokenStore: ts}
}

func (s *authService) Login(ctx core.Context, email string) (string, error) {
	// 1. 商業邏輯：先找使用者
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	// 2. 商業邏輯：生成 Token 並存入 Store (帶有計時器)
	token := "gen_uuid_here"
	userIDStr := strconv.FormatUint(uint64(user.ID), 10)

	err = s.tokenStore.Save(ctx, userIDStr, token, 24*time.Hour)

	return token, err
}

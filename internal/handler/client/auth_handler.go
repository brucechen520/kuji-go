package client

import (
	"net/http"

	clientS "github.com/brucechen520/kuji-go/internal/service/client"
	"github.com/brucechen520/kuji-go/internal/pkg/core"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *clientS.AuthService
}

func NewAuthHandler(as *clientS.AuthService) *AuthHandler {
	return &AuthHandler{authService: as}
}

func (h *AuthHandler) Login(c *gin.Context) {
	// 1. 解析參數 (這裡可以使用一個特定的 Request struct)
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. 呼叫 Service
	ctx := core.NewContext(c)
	defer core.ReleaseContext(ctx)

	token, err := h.authService.Login(ctx, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	// 3. 回傳結果
	c.JSON(http.StatusOK, gin.H{"token": token})
}

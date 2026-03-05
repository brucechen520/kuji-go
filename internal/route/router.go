package route

import (
	clientH "github.com/brucechen520/kuji-go/internal/handler/client"
	"github.com/gin-gonic/gin"
)

func NewRouter(authHandler *clientH.AuthHandler) *gin.Engine {
	r := gin.Default()

	v1ClientGroup := r.Group("/api/v1/client")
	{
		v1ClientGroup.POST("/login", authHandler.Login)
	}

	return r
}

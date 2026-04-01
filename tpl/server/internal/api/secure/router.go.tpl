package secure

import (
	"github.com/gin-gonic/gin"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
	"<xpfo{ .ModulePath }xpfo>/internal/controller/secure"
	"<xpfo{ .ModulePath }xpfo>/internal/middleware"
)

func Register(cfg *config.Config, router *gin.Engine) {
	group := router.Group("/api/v1/secure")
	group.Use(middleware.JWTAuth(cfg))
	group.GET("/ping", secure.Ping)
}

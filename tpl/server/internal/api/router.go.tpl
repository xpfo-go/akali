package api

import (
	"github.com/gin-gonic/gin"
	"<xpfo{ .ModulePath }xpfo>/internal/api/basic"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
	"<xpfo{ .ModulePath }xpfo>/internal/server"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := server.NewBasicRouter(cfg)

	// 注册路由
	basic.Register(cfg, router)

	return router
}

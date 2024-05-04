package api

import (
	"github.com/gin-gonic/gin"
	"<xpfo{ .ProjectName }xpfo>/internal/api/basic"
	"<xpfo{ .ProjectName }xpfo>/internal/config"
	"<xpfo{ .ProjectName }xpfo>/internal/server"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := server.NewBasicRouter(cfg)

	// 注册路由
	basic.Register(cfg, router)

	return router
}

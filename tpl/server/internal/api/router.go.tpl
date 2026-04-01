package api

import (
	"github.com/gin-gonic/gin"
	"<xpfo{ .ModulePath }xpfo>/internal/api/basic"
<xpfo{ if .EnableAuth }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/api/secure"
<xpfo{ end }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/config"
	"<xpfo{ .ModulePath }xpfo>/internal/server"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := server.NewBasicRouter(cfg)

	// 注册路由
	basic.Register(cfg, router)
<xpfo{ if .EnableAuth }xpfo>
	secure.Register(cfg, router)
<xpfo{ end }xpfo>

	return router
}

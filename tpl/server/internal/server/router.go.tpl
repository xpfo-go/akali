package server

import (
	"github.com/gin-gonic/gin"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
	"<xpfo{ .ModulePath }xpfo>/internal/middleware"
)

// NewRouterFunc ...
type NewRouterFunc func(cfg *config.Config) *gin.Engine

// NewBasicRouter ...
func NewBasicRouter(cfg *config.Config) *gin.Engine {
	if !cfg.Server.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	// disable console log color
	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// MW: request_id
	router.Use(middleware.RequestID())
<xpfo{ if .EnableRate }xpfo>
	router.Use(middleware.RateLimit(cfg))
<xpfo{ end }xpfo>

	return router
}

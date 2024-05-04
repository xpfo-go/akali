package server

import (
	"github.com/gin-gonic/gin"
	"<xpfo{ .ProjectName }xpfo>/internal/config"
	"<xpfo{ .ProjectName }xpfo>/internal/middleware"
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
	// MW: request_id
	router.Use(middleware.RequestID())

	return router
}

package api

import (
	"github.com/gin-gonic/gin"
	"{{ .ProjectName }}/config"
	"{{ .ProjectName }}/internal/api/basic"
	"{{ .ProjectName }}/pkg/server"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	router := server.NewBasicRouter(cfg)

	// 注册路由
	basic.Register(cfg, router)

	return router
}

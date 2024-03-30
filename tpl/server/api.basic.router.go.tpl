package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"net/http/pprof"

	"{{ .ProjectName }}/config"
	_ "{{ .ProjectName }}/docs"
	"{{ .ProjectName }}/internal/controller/basic"
)

func Register(cfg *config.Config, router *gin.Engine) {
	// basic
	router.GET("/ping", basic.Ping)
	router.GET("/health", basic.NewHealthHandleFunc(cfg))
	router.GET("/version", basic.Version)

	// metrics
	metricRouter := router.Group("/metrics")
	metricRouter.GET("", gin.WrapH(promhttp.Handler()))

	// pprof
	pprofRouter := router.Group("/debug/pprof")
	if !cfg.Server.IsDebug {
		pprofRouter.Use(gin.BasicAuth(gin.Accounts{
			"admin": "admin",
		}))
	}
	{
		pprofRouter.GET("/", pprofHandler(pprof.Index))
		pprofRouter.GET("/cmdline", pprofHandler(pprof.Cmdline))
		pprofRouter.GET("/profile", pprofHandler(pprof.Profile))
		pprofRouter.POST("/symbol", pprofHandler(pprof.Symbol))
		pprofRouter.GET("/symbol", pprofHandler(pprof.Symbol))
		pprofRouter.GET("/trace", pprofHandler(pprof.Trace))
		pprofRouter.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
		pprofRouter.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
		pprofRouter.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
		pprofRouter.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
		pprofRouter.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
		pprofRouter.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
	}

	// swagger docs
	if cfg.Server.IsDebug {
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

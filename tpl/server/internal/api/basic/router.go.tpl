package basic

import (
	"github.com/gin-gonic/gin"
<xpfo{ if .EnableMetrics }xpfo>
	"github.com/prometheus/client_golang/prometheus/promhttp"
<xpfo{ end }xpfo><xpfo{ if .EnableSwagger }xpfo>
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
<xpfo{ end }xpfo>
	"net/http"
	"net/http/pprof"
	"<xpfo{ .ModulePath }xpfo>/internal/config"

<xpfo{ if .EnableSwagger }xpfo>
	_ "<xpfo{ .ModulePath }xpfo>/docs"
<xpfo{ end }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/controller/basic"
)

func Register(cfg *config.Config, router *gin.Engine) {
	// basic
	router.GET("/ping", basic.Ping)
	router.GET("/health", basic.NewHealthHandleFunc(cfg))
	router.GET("/version", basic.Version)

<xpfo{ if .EnableMetrics }xpfo>
	// metrics
	metricRouter := router.Group("/metrics")
	metricRouter.GET("", gin.WrapH(promhttp.Handler()))
<xpfo{ end }xpfo>

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

<xpfo{ if .EnableSwagger }xpfo>
	// swagger docs
	if cfg.Server.IsDebug {
		url := ginSwagger.URL("/swagger/doc.json")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	}
<xpfo{ end }xpfo>
}

func pprofHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

package basic

import (
<xpfo{ if .EnableRedis }xpfo>
	"context"
<xpfo{ end }xpfo>
	"github.com/gin-gonic/gin"
	"net/http"
<xpfo{ if .EnableRedis }xpfo>
	"time"
<xpfo{ end }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/config"
<xpfo{ if .EnableMySQL }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/database"
<xpfo{ end }xpfo>
<xpfo{ if .EnableRedis }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/cache"
<xpfo{ end }xpfo>
)

// NewHealthHandleFunc godoc
// @Summary health for server health check
// @Description /health to make sure the server is health
// @ID health
// @Tags basic
// @Accept json
// @Produce json
// @Success 200 {string} string message
// @Failure 500 {string} string message
// @Header 200 {string} X-Request-Id "the request id"
// @Router /health [get]
func NewHealthHandleFunc(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
<xpfo{ if .EnableMySQL }xpfo>
		// check database
		if err := database.GetDefaultDBClient().TestConnection(); err != nil {
			c.String(http.StatusInternalServerError, "DefaultDBClient connect failed.")
			return
		}
<xpfo{ end }xpfo>
<xpfo{ if .EnableRedis }xpfo>
		if cache.DefaultRedisClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := cache.DefaultRedisClient.Ping(ctx).Err(); err != nil {
				c.String(http.StatusInternalServerError, "Redis connect failed.")
				return
			}
		}
<xpfo{ end }xpfo>

		c.String(http.StatusOK, "ok")
		return
	}
}

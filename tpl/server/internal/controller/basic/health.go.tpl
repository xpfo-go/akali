package basic

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"<xpfo{ .ProjectName }xpfo>/internal/config"
	"<xpfo{ .ProjectName }xpfo>/internal/database"
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
		// check database
		if err := database.GetDefaultDBClient().TestConnection(); err != nil {
			c.String(http.StatusInternalServerError, "DefaultDBClient connect failed.")
			return
		}

		c.String(http.StatusOK, "ok")
		return
	}
}

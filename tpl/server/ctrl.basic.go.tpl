package basic

import (
	"github.com/gin-gonic/gin"
	"os"
	"{{ .ProjectName }}/internal/version"
	"time"
)

// Ping godoc
// @Summary ping-pong for alive test
// @Description /ping to get response from {{ .ProjectName }}, make sure the server is alive
// @ID ping
// @Tags basic
// @Accept json
// @Produce json
// @Success 200
// @Header 200 {string} X-Request-Id "the request id"
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// Version godoc
// @Summary version for identify
// @Description /version to get the version of {{ .ProjectName }}
// @ID version
// @Tags basic
// @Accept json
// @Produce json
// @Success 200
// @Header 200 {string} X-Request-Id "the request id"
// @Router /version [get]
func Version(c *gin.Context) {
	runEnv := os.Getenv("RUN_ENV")
	now := time.Now()
	c.JSON(200, gin.H{
		"version":   version.Version,
		"commit":    version.Commit,
		"buildTime": version.BuildTime,
		"goVersion": version.GoVersion,
		"env":       runEnv,
		// return the date and timestamp
		"timestamp": now.Unix(),
		"date":      now,
	})
}

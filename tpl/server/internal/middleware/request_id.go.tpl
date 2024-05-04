package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/xpfo-go/logs"
	"<xpfo{ .ProjectName }xpfo>/internal/util"
)

// RequestID add the request_id for each api request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		logs.Debug("Middleware: RequestID")

		requestID := c.GetHeader(util.RequestIDHeaderKey)
		if requestID == "" || len(requestID) != 32 {
			requestID = util.GenUUID4()
		}
		util.SetRequestID(c, requestID)
		c.Writer.Header().Set(util.RequestIDHeaderKey, requestID)

		c.Next()
	}
}

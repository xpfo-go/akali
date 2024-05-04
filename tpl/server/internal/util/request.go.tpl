package util

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

// ErrNilRequestBody ...
var ErrNilRequestBody = errors.New("request Body is nil")

// ReadRequestBody will return the body in []byte, without change the origin body
func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, ErrNilRequestBody
	}

	body, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(body))
	return body, err
}

// GetRequestID ...
func GetRequestID(c *gin.Context) string {
	return c.GetString(RequestIDKey)
}

// SetRequestID ...
func SetRequestID(c *gin.Context, requestID string) {
	c.Set(RequestIDKey, requestID)
}

// GetClientID ...
func GetClientID(c *gin.Context) string {
	return c.GetString(ClientIDKey)
}

// SetClientID ...
func SetClientID(c *gin.Context, clientID string) {
	c.Set(ClientIDKey, clientID)
}

// GetError ...
func GetError(c *gin.Context) (interface{}, bool) {
	return c.Get(ErrorIDKey)
}

// SetError ...
func SetError(c *gin.Context, err error) {
	c.Set(ErrorIDKey, err)
}

// BasicAuthAuthorizationHeader ...
func BasicAuthAuthorizationHeader(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString(StringToBytes(base))
}

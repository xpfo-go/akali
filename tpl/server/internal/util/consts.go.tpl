package util

// RequestIDKey ...
const (
	RequestIDKey       = "request_id"
	RequestIDHeaderKey = "X-Request-Id"

	ClientIDKey = "client_id"

	ErrorIDKey = "err"

	// NeverExpiresUnixTime 永久有效期，使用2100.01.01 00:00:00 的unix time作为永久有效期的表示，单位秒
	// time.Date(2100, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	NeverExpiresUnixTime = 4102444800
)

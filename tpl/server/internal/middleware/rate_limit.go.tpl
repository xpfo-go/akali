package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors   = map[string]*visitor{}
	visitorsMu sync.Mutex
)

func getVisitor(ip string, r rate.Limit, burst int) *rate.Limiter {
	visitorsMu.Lock()
	defer visitorsMu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(r, burst)
		visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func cleanupVisitors() {
	for {
		time.Sleep(1 * time.Minute)
		visitorsMu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 5*time.Minute {
				delete(visitors, ip)
			}
		}
		visitorsMu.Unlock()
	}
}

func RateLimit(cfg *config.Config) gin.HandlerFunc {
	go cleanupVisitors()

	if !cfg.RateLimit.Enabled {
		return func(c *gin.Context) { c.Next() }
	}

	rps := cfg.RateLimit.RPS
	if rps <= 0 {
		rps = 50
	}
	burst := cfg.RateLimit.Burst
	if burst <= 0 {
		burst = 100
	}

	limit := rate.Limit(rps)
	return func(c *gin.Context) {
		limiter := getVisitor(c.ClientIP(), limit, burst)
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

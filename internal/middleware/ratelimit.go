package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	limiter  *time.Ticker
	lastSeen time.Time
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func RateLimitMiddleware(ctx context.Context, limit time.Duration) gin.HandlerFunc {
	go cleanupVisitors(ctx)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			v = &visitor{limiter: time.NewTicker(limit), lastSeen: time.Now()}
			visitors[ip] = v
		}
		v.lastSeen = time.Now()
		mu.Unlock()

		select {
		case <-v.limiter.C:
			c.Next()
		default:
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}
	}
}

func cleanupVisitors(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					v.limiter.Stop()
					delete(visitors, ip)
				}
			}
			mu.Unlock()

		}
	}
}

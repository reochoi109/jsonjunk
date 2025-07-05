package middleware

import (
	"context"
	"jsonjunk/internal/model"
	"jsonjunk/pkg/idgen"

	"github.com/gin-gonic/gin"
)

func UseTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := idgen.GenerateTraceID()
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "no-referrer")
		c.Next()
	}
}

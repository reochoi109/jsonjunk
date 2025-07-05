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

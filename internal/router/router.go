package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "jsonjunk/docs"
	"jsonjunk/internal/handler"
	"jsonjunk/internal/helper"
	"jsonjunk/internal/middleware"
	"jsonjunk/internal/model"
	"jsonjunk/internal/service"
)

func NewRouter(mode string, svc service.PasteService) *gin.Engine {
	gin.SetMode(mode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.SecurityHeaders())

	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		log := model.WithContext(c.Request.Context())
		log.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.UseTraceID())

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterAPI(r.Group("/raw"), svc)
	api := r.Group("/api/v1")
	RegisterPastes(api, svc)

	return r
}

func Run(ctx context.Context, mode string, svc service.PasteService) {
	r := NewRouter(mode, svc)

	srv := &http.Server{
		Addr:         ":8080", // config에서 주입
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	helper.GracefulShutdown(ctx, srv)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}

	// HTTPs
	// if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
	// 	panic(err)
	// }
}

func RegisterPastes(api *gin.RouterGroup, svc service.PasteService) {
	paste := api.Group("/paste")
	paste.GET("/list", handler.GetSearchPastedList(svc))
	paste.GET("/type", handler.GetExpireType)
	paste.GET("/:id", handler.GetPasteHandler(svc))
	paste.PUT("/:id", handler.UpdatePasteHandler(svc))
	paste.POST("", handler.CreatePasteHandler(svc))
}

func RegisterAPI(api *gin.RouterGroup, svc service.PasteService) {
	api.GET("/:id", handler.GetPasteJsonHandler(svc))
}

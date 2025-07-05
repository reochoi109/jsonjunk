package router

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "jsonjunk/docs"
	"jsonjunk/internal/handler"
	"jsonjunk/internal/service"
	logger "jsonjunk/pkg/logging"
)

func Run(ctx context.Context, svc service.PasteService) {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/raw")
	RegisterAPI(api, svc)

	group := r.Group("/api/v1")
	RegisterPastes(group, svc)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Log.Error("server shutdown failed", zap.Error(err))
		} else {
			logger.Log.Info("server shutdown completed")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
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

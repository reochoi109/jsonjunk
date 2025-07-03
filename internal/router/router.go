package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "jsonjunk/docs"
	"jsonjunk/internal/handler"
	"jsonjunk/internal/service"
)

func Run(svc service.PasteService) {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	RegisterPastes(api, svc)
	r.Run(":8080")
}

func RegisterPastes(api *gin.RouterGroup, svc service.PasteService) {
	paste := api.Group("/paste")
	paste.GET("/list", handler.GetSearchPastedList(svc))
	paste.GET("/type", handler.GetExpireType)
	paste.GET("/:id", handler.GetPasteHandler(svc))
	paste.PUT("/:id", handler.UpdatePasteHandler(svc))
	paste.POST("", handler.CreatePasteHandler(svc))
}

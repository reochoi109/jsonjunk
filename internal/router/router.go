package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "jsonjunk/docs"
	"jsonjunk/internal/handler"
)

func Run() {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	RegisterPastes(api)
	r.Run(":8080")
}

func RegisterPastes(api *gin.RouterGroup) {
	paste := api.Group("/paste")
	paste.GET("/:id", handler.HandleGetPaste) // 등록된 paste 조회
	paste.POST("", handler.HandleCreatePaste) // 신규 paste 조회
}

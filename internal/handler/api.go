package handler

import (
	"context"
	"jsonjunk/internal/model"
	"jsonjunk/internal/service"
	"jsonjunk/pkg/idgen"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaste godoc
//
//	@Summary		데이터를 조회한다.
//	@Description	데이터를 조회한다.
//	@Tags			api
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Paste ID"
//	@Router			/raw/{id} [get]
func GetPasteJsonHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextRequestID, idgen.GenerateUUID())
		id := c.Param("id")
		paste, err := svc.GetPasteByID(ctx, id)
		if err != nil || paste == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		response := model.PasteResponse{
			ID:        paste.ID,
			Title:     paste.Title,
			Language:  paste.Language,
			CreatedAt: paste.CreatedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt: paste.ExpiresAt.Format("2006-01-02 15:04:05"),
			Content:   paste.Content,
		}
		c.JSON(http.StatusOK, response)
	}
}

package handler

import (
	"context"
	"jsonjunk/internal/model"
	"jsonjunk/internal/service"
	"jsonjunk/pkg/idgen"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPasteRaw godoc
//
//	@Summary		원본 Paste 콘텐츠 조회
//	@Description	줄바꿈과 포맷 그대로의 원본 텍스트를 반환합니다.
//	@Tags			api
//	@Produce		plain
//	@Param			id	path	string	true	"Paste ID"
//	@Success		200		{string}	string	"원본 콘텐츠"
//	@Failure		404		{object}	map[string]string
//	@Router			/raw/{id} [get]
func GetPasteJsonHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, idgen.GenerateTraceID())
		log := model.WithContext(ctx)
		log.Debug("GetPasteJsonHandler called [Start]")
		defer log.Debug("GetPasteJsonHandler [End]")

		id := c.Param("id")
		paste, err := svc.GetPasteByID(ctx, id)
		if err != nil || paste == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
			return
		}
		c.String(http.StatusOK, paste.Content)
	}
}

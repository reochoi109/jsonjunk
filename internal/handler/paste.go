package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaste godoc
//	@Summary		Paste 조회
//	@Description	ID를 통해 저장된 Paste 내용을 조회합니다.
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Paste ID"
//	@Success		200	{object}	model.PasteResponse
//	@Failure		404	{object}	model.ErrorResponse
//	@Router			/paste/{id} [get]
func HandleGetPaste(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// CreatePaste godoc
//	@Summary		Paste 생성
//	@Description	새로운 Paste 텍스트를 생성하고 저장합니다.
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			paste	body		model.PasteRequest	true	"Paste Content"
//	@Success		200		{object}	model.PasteResponse
//	@Failure		400		{object}	model.ErrorResponse
//	@Router			/paste [post]
func HandleCreatePaste(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

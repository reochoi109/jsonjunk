package handler

import (
	"jsonjunk/internal/model"
	"jsonjunk/internal/service"
	"jsonjunk/pkg/idgen"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPaste godoc
//
//	@Summary		Paste expire type 조회
//	@Description	Paste expire type 조회합니다.
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	model.ResponseFormat{data=[]model.ExpiredTypeResponse}
//	@Failure		404	{object}	model.ResponseFormat
//	@Router			/paste/type [get]
func GetExpireType(c *gin.Context) {
	types := []model.ExpiredTypeResponse{
		{
			Type: int(model.Expire6Hours),
			Name: "6h",
		},
		{
			Type: int(model.Expire12Hours),
			Name: "12h",
		},
		{
			Type: int(model.Expire1Day),
			Name: "1day",
		},
		{
			Type: int(model.Expire7Days),
			Name: "1week",
		},
	}
	c.JSON(http.StatusOK, types)
}

// GetPaste godoc
//
//	@Summary		Paste 조회
//	@Description	ID를 통해 저장된 Paste 내용을 조회합니다.
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Paste ID"
//	@Success		200	{object}	model.PasteResponse
//	@Failure		404	{object}	model.PastErrorResponse
//	@Router			/paste/{id} [get]
func GetPasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		paste, err := svc.GetPaste(id)
		if err != nil || paste == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
			return
		}
		c.JSON(http.StatusOK, paste)
	}
}

// CreatePaste godoc
//
//	@Summary		Paste 생성
//	@Description	새로운 Paste 텍스트를 생성하고 저장합니다.
//	@Description	expire
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			paste	body		model.PasteRequest	true	"Paste Content"
//	@Success		200		{object}	model.ResponseFormat
//	@Failure		400		{object}	model.ResponseFormat
//	@Router			/paste [post]
func CreatePasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.PasteRequest
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}

		if err := svc.CreatePaste(model.Paste{
			ID:        idgen.GenerateUUID(),
			Content:   req.Content,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(req.Expire.Duration()),
		}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Created"})
	}
}

// CreatePaste godoc
//
//	@Summary		test
//	@Description	test
//	@Tags			test
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	model.ResponseFormat
//	@Failure		400		{object}	model.ResponseFormat
//	@Router			/paste/test/list [get]
func TestSearchPastedAll(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := svc.TestSearch()
		c.JSON(http.StatusOK, data)
	}
}

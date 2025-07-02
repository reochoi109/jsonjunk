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
//	@Tags			pastes:type
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
	model.HandleResponse(c, http.StatusOK, model.Success, types)
}

// CreatePaste godoc
//
//	@Summary		Paste 목록 조회
//	@Description	Paste 목록 조회 요청
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	model.ResponseFormat{data=[]model.PasteResponse}
//	@Failure		400		{object}	model.ResponseFormat
//	@Router			/paste/list [get]
func GetSearchPastedList(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		datas, _ := svc.SearchPasteList()
		response := make([]model.PasteResponse, len(datas))
		for i, v := range datas {
			response[i] = model.PasteResponse{
				ID:        v.ID,
				Title:     v.Title,
				CreatedAt: v.CreatedAt.Format("2006-01-02 15:04:05"),
				ExpiresAt: v.ExpiresAt.Format("2006-01-02 15:04:05"),
			}
		}
		model.HandleResponse(c, http.StatusOK, model.Success, response)
	}
}

// GetPaste godoc
//
//	@Summary		Paste 조회
//	@Description	ID를 통해 저장된 Paste 내용을 조회합니다.
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Paste ID"
//	@Success		200	{object}	model.ResponseFormat{data=model.PasteResponse}
//	@Failure		404	{object}	model.ResponseFormat
//	@Router			/paste/{id} [get]
func GetPasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		paste, err := svc.SearchPaste(id)
		if err != nil || paste == nil {
			model.HandleResponse(c, http.StatusNotFound, model.ErrorPasteNotFound, nil)
			return
		}
		response := model.PasteResponse{
			ID:        paste.ID,
			Title:     paste.Title,
			CreatedAt: paste.CreatedAt.Format("2006-01-02 15:04:05"),
			ExpiresAt: paste.ExpiresAt.Format("2006-01-02 15:04:05"),
			Content:   paste.Content,
		}
		model.HandleResponse(c, http.StatusOK, model.Success, response)
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
		if err := c.ShouldBindJSON(&req); err != nil {
			model.HandleResponse(c, http.StatusBadRequest, model.ErrorValidationFailed, nil)
			return
		}

		if err := svc.CreatePaste(model.Paste{
			ID:        idgen.GenerateUUID(),
			Title:     req.Title,
			Content:   req.Content,
			CreatedAt: time.Now().UTC(),
			ExpiresAt: time.Now().UTC().Add(req.Expire.Duration()),
		}); err != nil {
			model.HandleResponse(c, http.StatusBadRequest, model.ErrorPasteCreateFailed, nil)
			return
		}
		model.HandleResponse(c, http.StatusCreated, model.Success, nil)
	}
}

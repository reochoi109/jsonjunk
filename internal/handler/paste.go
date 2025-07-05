package handler

import (
	"context"
	"errors"
	"fmt"
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
//	@Router			/api/v1/paste/type [get]
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
//	@Router			/api/v1/paste/list [get]
func GetSearchPastedList(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, idgen.GenerateUUID())
		datas, err := svc.GetListPastes(ctx)
		if err != nil {
			if errors.Is(err, model.ErrDatabase) {
				model.HandleResponse(c, http.StatusInternalServerError, model.ErrorDatabase, nil)
				return
			}
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorInternalServer, nil)
			return
		}

		response := make([]model.PasteResponse, len(datas))
		for i, v := range datas {
			response[i] = model.PasteResponse{
				ID:        v.ID,
				Title:     v.Title,
				Language:  v.Language,
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
//	@Router			/api/v1/paste/{id} [get]
func GetPasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, idgen.GenerateUUID())
		id := c.Param("id")
		paste, err := svc.GetPasteByID(ctx, id)
		if err != nil || paste == nil {
			model.HandleResponse(c, http.StatusNotFound, model.ErrorPasteNotFound, nil)
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
		model.HandleResponse(c, http.StatusOK, model.Success, response)
	}
}

// CreatePaste godoc
//
//	@Summary		Paste 생성
//	@Description	새로운 Paste 텍스트를 생성하고 저장합니다.
//	@Tags			pastes
//	@Accept			json
//	@Produce		json
//	@Param			paste	body		model.PasteRequest	true	"Paste Content"
//	@Success		200		{object}	model.ResponseFormat
//	@Failure		400		{object}	model.ResponseFormat
//	@Router			/api/v1/paste [post]
func CreatePasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, idgen.GenerateUUID())

		var req model.PasteRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			model.HandleResponse(c, http.StatusBadRequest, model.ErrorValidationFailed, nil)
			return
		}

		paste := model.Paste{
			ID:                 idgen.GenerateUUID(),
			Title:              req.Title,
			Language:           req.Language,
			Content:            req.Content,
			CreatedAt:          time.Now().UTC(),
			AnonymousExpiresAt: time.Now().UTC().Add(req.Expire.Duration()),
		}

		err := svc.RegisterPaste(ctx, paste)
		switch {
		case errors.Is(err, model.ErrDuplicatePasteID):
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorDatabase, nil)
		case errors.Is(err, model.ErrInsertFailed):
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorDatabase, nil)
		case err != nil:
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorInternalServer, nil)
		default:
			model.HandleResponse(c, http.StatusCreated, model.SuccessPasteCreated, nil)
		}
	}
}

// UpdatePasteHandler godoc
//
//	@Summary		Paste 업데이트 : 테스트 용
//	@Description	Paste 텍스트를 업데이트 및 저장합니다.
//	@Tags			pastes:test
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Paste ID"
//	@Param			paste	body		model.PasteUpdateRequest	true	"Paste update Content"
//	@Success		200		{object}	model.ResponseFormat
//	@Failure		400		{object}	model.ResponseFormat
//	@Router			/api/v1/paste/{id} [put]
func UpdatePasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, idgen.GenerateUUID())

		id := c.Param("id")
		var req model.PasteUpdateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			fmt.Println(err)
			model.HandleResponse(c, http.StatusBadRequest, model.ErrorValidationFailed, nil)
			return
		}
		fields := make(map[string]interface{})

		if req.Title != nil {
			fields["title"] = *req.Title
		}

		if req.Content != nil {
			fields["content"] = *req.Content
		}

		if req.Language != nil {
			fields["language"] = *req.Language
		}

		if req.Expire != nil {
			if !req.Expire.IsValid() {
				model.HandleResponse(c, http.StatusBadRequest, model.ErrorValidationFailed, nil)
				return
			}
			fields["expires_at"] = time.Now().Add(req.Expire.Duration())
		}

		if len(fields) == 0 {
			model.HandleResponse(c, http.StatusBadRequest, model.ErrorNoUpdatableField, nil)
			return
		}

		updated, err := svc.UpdatePasteByID(ctx, id, fields)
		switch {
		case errors.Is(err, model.ErrPasteNotFound):
			model.HandleResponse(c, http.StatusNotFound, model.ErrorPasteNotFound, nil)
		case errors.Is(err, model.ErrDatabase):
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorDatabase, nil)
		case err != nil:
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorInternalServer, nil)
		default:
			model.HandleResponse(c, http.StatusOK, model.SuccessPasteUpdated, updated)
		}
	}
}

// RemovePasteHandler godoc
//
//	@Summary		Paste 삭제 : : 테스트 용
//	@Description	Paste 텍스트를 삭제합니다.
//	@Tags			pastes:test
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string	true	"Paste ID"
//	@Success		200		{object}	model.ResponseFormat
//	@Failure		400		{object}	model.ResponseFormat
//	@Router			/api/v1/paste/{id} [delete]
func RemovePasteHandler(svc service.PasteService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), model.ContextTraceID, idgen.GenerateUUID())
		id := c.Param("id")
		err := svc.RemovePasteByID(ctx, id)
		switch {
		case errors.Is(err, model.ErrPasteNotFound):
			model.HandleResponse(c, http.StatusNotFound, model.ErrorPasteNotFound, nil)
		case err != nil:
			model.HandleResponse(c, http.StatusInternalServerError, model.ErrorDatabase, nil)
		default:
			model.HandleResponse(c, http.StatusNoContent, model.Success, nil)
		}
	}
}

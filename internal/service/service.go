package service

import (
	"context"
	"jsonjunk/internal/model"
	"jsonjunk/internal/repository"
	"time"

	"go.uber.org/zap"
)

const mongoTimeout = 5 * time.Second

type PasteService interface {
	RegisterPaste(ctx context.Context, req model.PasteRequest) error
	GetPasteByID(ctx context.Context, id string) (*model.PasteResponse, error)
	GetListPastes(ctx context.Context) ([]model.PasteResponse, error)
	UpdatePasteByID(ctx context.Context, id string, fields map[string]interface{}) (paste model.Paste, err error)
	RemovePasteByID(ctx context.Context, id string) error
}

type pasteService struct {
	repo repository.Repository
}

func NewPasteService(ctx context.Context, repo repository.Repository) PasteService {
	service := &pasteService{repo: repo}
	service.internal(ctx)
	return service
}

func (s *pasteService) RegisterPaste(ctx context.Context, req model.PasteRequest) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	paste := model.NewPasteFromRequest(req)

	log := model.WithContext(ctx)
	log.Debug("Creating new paste [Start]", zap.String("id", paste.ID))
	defer log.Debug("Creating new paste [End]", zap.String("id", paste.ID))

	if err := s.repo.InsertPaste(ctx, paste); err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to insert paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", paste.ID), zap.Error(err))
		return err
	}

	return nil
}

func (s *pasteService) GetPasteByID(ctx context.Context, id string) (*model.PasteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("Searching paste by ID [Start]", zap.String("id", id))

	paste, err := s.repo.SearchPasteByID(ctx, id)
	if err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to search paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", id), zap.Error(err))
		return nil, err
	}
	response := model.NewPasteResponse(*paste)
	log.Debug("Searching paste by ID [End]", zap.String("id", id))
	return &response, nil
}

func (s *pasteService) GetListPastes(ctx context.Context) ([]model.PasteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("GetListPastes called [Start]")
	defer log.Debug("GetListPastes [End]")

	pastes, err := s.repo.SearchPasteList(ctx)
	if err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to search paste list", zap.String(string(model.ContextTraceID), traceID.(string)), zap.Error(err))
		return nil, err
	}

	response := make([]model.PasteResponse, len(pastes))
	for i, p := range pastes {
		response[i] = model.NewPasteListResponse(*p)
	}
	return response, nil
}

func (s *pasteService) UpdatePasteByID(ctx context.Context, id string, fields map[string]interface{}) (model.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("Modifying paste [Start]", zap.String("id", id))
	defer log.Debug("Modifying paste [End]", zap.String("id", id))

	paste, err := s.repo.UpdatePasteByID(ctx, id, fields)
	if err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to modify paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", id), zap.Error(err))
		return paste, err
	}
	return paste, nil
}

func (s *pasteService) RemovePasteByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("Deleting paste [Start]", zap.String("id", id))
	defer log.Debug("Deleting paste [End]", zap.String("id", id))

	if err := s.repo.DeletePasteByID(ctx, id); err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to delete paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", id), zap.Error(err))
		return err
	}
	return nil
}

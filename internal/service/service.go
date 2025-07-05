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
	RegisterPaste(ctx context.Context, p model.Paste) error
	GetPasteByID(ctx context.Context, id string) (*model.Paste, error)
	GetListPastes(ctx context.Context) ([]*model.Paste, error)
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

func (s *pasteService) RegisterPaste(ctx context.Context, p model.Paste) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("Creating new paste [Start]", zap.String("id", p.ID))

	if err := s.repo.InsertPaste(ctx, p); err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to insert paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", p.ID), zap.Error(err))
		return err
	}
	log.Debug("Creating new paste [End]", zap.String("id", p.ID))
	return nil
}

func (s *pasteService) GetPasteByID(ctx context.Context, id string) (*model.Paste, error) {
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
	log.Debug("Searching paste by ID [End]", zap.String("id", id))
	return paste, nil
}

func (s *pasteService) GetListPastes(ctx context.Context) ([]*model.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()
	log := model.WithContext(ctx)
	log.Debug("GetListPastes called [Start]")

	pastes, err := s.repo.SearchPasteList(ctx)
	if err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to search paste list", zap.String(string(model.ContextTraceID), traceID.(string)), zap.Error(err))
		return nil, err
	}
	log.Debug("GetListPastes [End]")
	return pastes, nil
}

func (s *pasteService) UpdatePasteByID(ctx context.Context, id string, fields map[string]interface{}) (model.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("Modifying paste [Start]", zap.String("id", id))

	paste, err := s.repo.UpdatePasteByID(ctx, id, fields)
	if err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to modify paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", id), zap.Error(err))
		return paste, err
	}
	log.Debug("Modifying paste [End]", zap.String("id", id))
	return paste, nil
}

func (s *pasteService) RemovePasteByID(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Debug("Deleting paste [Start]", zap.String("id", id))

	if err := s.repo.DeletePasteByID(ctx, id); err != nil {
		traceID := ctx.Value(model.ContextTraceID)
		log.Error("Failed to delete paste", zap.String(string(model.ContextTraceID), traceID.(string)), zap.String("id", id), zap.Error(err))
		return err
	}
	log.Debug("Deleting paste [End]", zap.String("id", id))
	return nil
}

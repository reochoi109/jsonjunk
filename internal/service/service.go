package service

import (
	"context"
	"fmt"
	"jsonjunk/internal/model"
	"jsonjunk/internal/repository"
	"time"

	"go.uber.org/zap"
)

const mongoTimeout = 5 * time.Second

type PasteService interface {
	CreatePaste(ctx context.Context, p model.Paste) error
	SearchPaste(ctx context.Context, id string) (*model.Paste, error)
	SearchPasteList(ctx context.Context) ([]*model.Paste, error)
	ModifyPaste(ctx context.Context, id string, fields map[string]interface{}) (paste model.Paste, err error)
	DeletePaste(ctx context.Context, id string) error
}

type pasteService struct {
	repo repository.Repository
}

func NewPasteService(repo repository.Repository) PasteService {
	return &pasteService{repo: repo}
}

func (s *pasteService) CreatePaste(ctx context.Context, p model.Paste) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Info("Creating new paste", zap.String("id", p.ID))

	if err := s.repo.Insert(ctx, p); err != nil {
		log.Error("Failed to insert paste", zap.String("id", p.ID), zap.Error(err))
		return fmt.Errorf("create paste failed: %w", err)
	}
	return nil
}

func (s *pasteService) SearchPaste(ctx context.Context, id string) (*model.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Info("Searching paste by ID", zap.String("id", id))

	paste, err := s.repo.SearchPasteByID(ctx, id)
	if err != nil {
		log.Error("Failed to search paste", zap.String("id", id), zap.Error(err))
		return nil, fmt.Errorf("search paste failed: id=%s: %w", id, err)
	}
	return paste, nil
}

func (s *pasteService) SearchPasteList(ctx context.Context) ([]*model.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()
	log := model.WithContext(ctx)
	log.Info("Searching all pastes")

	pastes, err := s.repo.SearchPasteList(ctx)
	if err != nil {
		log.Error("Failed to search paste list", zap.Error(err))
		return nil, fmt.Errorf("search paste list failed: %w", err)
	}
	return pastes, nil
}

func (s *pasteService) ModifyPaste(ctx context.Context, id string, fields map[string]interface{}) (model.Paste, error) {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Info("Modifying paste", zap.String("id", id))

	paste, err := s.repo.ModifyPaste(ctx, id, fields)
	if err != nil {
		log.Error("Failed to modify paste", zap.String("id", id), zap.Error(err))
		return paste, fmt.Errorf("modify paste failed: id=%s: %w", id, err)
	}
	return paste, nil
}

func (s *pasteService) DeletePaste(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout)
	defer cancel()

	log := model.WithContext(ctx)
	log.Info("Deleting paste", zap.String("id", id))

	if err := s.repo.DeletePaste(ctx, id); err != nil {
		log.Error("Failed to delete paste", zap.String("id", id), zap.Error(err))
		return fmt.Errorf("delete paste failed: id=%s: %w", id, err)
	}
	return nil
}

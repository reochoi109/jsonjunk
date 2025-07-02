package service

import (
	"jsonjunk/internal/model"
	"jsonjunk/internal/repository"
)

type PasteService interface {
	CreatePaste(p model.Paste) error
	GetPaste(id string) (*model.Paste, error)
	TestSearch() ([]*model.Paste, error)
}

type pasteService struct {
	repo repository.Repository
}

func NewPasteService(repo repository.Repository) PasteService {
	return &pasteService{repo: repo}
}

func (s *pasteService) CreatePaste(p model.Paste) error {
	return s.repo.Insert(p)
}

func (s *pasteService) GetPaste(id string) (*model.Paste, error) {
	return s.repo.SearchPasteByID(id)
}

func (s *pasteService) TestSearch() ([]*model.Paste, error) {
	return s.repo.TestSearchPastedAll()
}

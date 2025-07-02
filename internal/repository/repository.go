package repository

import "jsonjunk/internal/model"

type Repository interface {
	Insert(p model.Paste) error
	SearchPasteByID(id string) (*model.Paste, error)
	SearchPasteList() ([]*model.Paste, error)
}

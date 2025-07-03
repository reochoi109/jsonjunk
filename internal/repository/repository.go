package repository

import (
	"context"
	"jsonjunk/internal/model"
)

type Repository interface {
	InsertPaste(ctx context.Context, p model.Paste) error
	SearchPasteByID(ctx context.Context, id string) (*model.Paste, error)
	SearchPasteList(ctx context.Context) ([]*model.Paste, error)
	UpdatePasteByID(ctx context.Context, id string, fields map[string]interface{}) (model.Paste, error)
	DeletePasteByID(ctx context.Context, id string) error
}

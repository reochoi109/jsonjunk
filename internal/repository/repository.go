package repository

import (
	"context"
	"jsonjunk/internal/model"
)

type Repository interface {
	Insert(ctx context.Context, p model.Paste) error
	SearchPasteByID(ctx context.Context, id string) (*model.Paste, error)
	SearchPasteList(ctx context.Context) ([]*model.Paste, error)
	ModifyPaste(ctx context.Context, id string, fields map[string]interface{}) (paste model.Paste, err error)
	DeletePaste(ctx context.Context, id string) error
}

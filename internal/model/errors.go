package model

import "errors"

var (
	ErrPasteNotFound    = errors.New("paste not found")
	ErrDatabase         = errors.New("database error")
	ErrDuplicatePasteID = errors.New("duplicate paste ID")
	ErrInsertFailed     = errors.New("failed to insert paste")
)

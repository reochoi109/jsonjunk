package model

import (
	"jsonjunk/pkg/idgen"
	"time"
)

type Paste struct {
	MongoID            string    `bson:"_id,omitempty" json:"mongo_id,omitempty"`               // Mongo의 내부 ID
	ID                 string    `bson:"id" json:"id"`                                          // ID
	Title              string    `bson:"title" json:"title"`                                    // 제목
	Language           string    `bson:"language,omitempty" json:"language" binding:"required"` // 언어
	Content            string    `bson:"content" json:"content"`                                // content
	UserID             string    `bson:"user_id,omitempty" json:"user_id,omitempty"`            // 사용자 ID (선택)
	CreatedAt          time.Time `bson:"created_at" json:"created_at"`                          // 생성일
	ExpiresAt          time.Time `bson:"expires_at,omitempty" json:"expires_at"`                // 사용자 만료일
	AnonymousExpiresAt time.Time `bson:"anonymous_expires_at,omitempty"`                        // 익명 사용자 만료일
}

func NewPasteFromRequest(req PasteRequest) Paste {
	// TODO :익명 사용자
	// TODO :가입 사용자
	now := time.Now().UTC()
	return Paste{
		ID:                 idgen.GeneratePasteID(),
		Title:              req.Title,
		Language:           req.Language,
		Content:            req.Content,
		CreatedAt:          now,
		AnonymousExpiresAt: now.Add(req.Expire.Duration()),
	}
}

func NewPasteResponse(p Paste) PasteResponse {
	return PasteResponse{
		ID:        p.ID,
		Title:     p.Title,
		Language:  p.Language,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		ExpiresAt: p.ExpiresAt.Format("2006-01-02 15:04:05"),
		Content:   p.Content,
	}
}

func NewPasteListResponse(p Paste) PasteResponse {
	return PasteResponse{
		ID:        p.ID,
		Title:     p.Title,
		Language:  p.Language,
		CreatedAt: p.CreatedAt.Format("2006-01-02 15:04:05"),
		ExpiresAt: p.ExpiresAt.Format("2006-01-02 15:04:05"),
	}
}

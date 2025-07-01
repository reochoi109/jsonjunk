package model

import "time"

type Paste struct {
	MongoID   string    `bson:"_id,omitempty"`        // Mongo의 내부 ID
	ID        string    `bson:"id"`                   // ID
	Content   string    `bson:"content"`              // content
	UserID    string    `bson:"user_id,omitempty"`    // 선택: 나중에 사용자 연결 시
	CreatedAt time.Time `bson:"created_at"`           // 생성 일자
	ExpiresAt time.Time `bson:"expires_at,omitempty"` // 유효기간
}

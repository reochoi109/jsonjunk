package model

import "time"

type Paste struct {
	MongoID   string    `bson:"_id,omitempty" json:"mongo_id,omitempty"`    // Mongo의 내부 ID
	ID        string    `bson:"id" json:"id"`                               // ID
	Title     string    `bson:"title" json:"title"`                         // 제목
	Content   string    `bson:"content" json:"content"`                     // content
	UserID    string    `bson:"user_id,omitempty" json:"user_id,omitempty"` // 사용자 ID (선택)
	CreatedAt time.Time `bson:"created_at" json:"created_at"`               // 생성일
	ExpiresAt time.Time `bson:"expires_at,omitempty" json:"expires_at"`     // 만료일
}

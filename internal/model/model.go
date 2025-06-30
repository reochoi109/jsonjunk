package model

type PasteRequest struct {
	Content string `json:"content" example:"text"`
	Expire  string `json:"expire" enums:"6h,12h,1d,7d,forever" example:"6h" description:"만료 옵션 (6시간, 12시간, 1일, 7일, 영구저장)"`
}

type PasteResponse struct {
	ID        string `json:"id" example:"abc123"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message" example:"해당 Paste를 찾을 수 없습니다."`
}

package model

type ExpiredTypeResponse struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

type PasteRequest struct {
	Content string       `json:"content" binding:"required"`
	Expire  ExpireOption `json:"expire" binding:"required"`
}

type PasteResponse struct {
	ID        string `json:"id,omitempty" example:"abc123"`
	CreatedAt string `json:"created_at,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Content   string `json:"content,omitempty"`
}

type PastErrorResponse struct {
	Error string `json:"error"`
}

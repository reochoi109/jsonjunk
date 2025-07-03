package model

type ExpiredTypeResponse struct {
	Type int    `json:"type"`
	Name string `json:"name"`
}

type PasteRequest struct {
	Title    string       `json:"title" binding:"required"`
	Content  string       `json:"content" binding:"required"`
	Language string       `json:"language" binding:"required"`
	Expire   ExpireOption `json:"expire" binding:"required,oneof=1 2 3 4" enums:"1,2,3,4"`
}

type PasteUpdateRequest struct {
	Title    *string       `json:"title,omitempty"`
	Content  *string       `json:"content,omitempty"`
	Language *string       `json:"language"`
	Expire   *ExpireOption `json:"expire,omitempty"`
}
type PasteResponse struct {
	ID        string `json:"id,omitempty" example:"abc123"`
	Title     string `json:"title" binding:"required"`
	Language  string `json:"language" binding:"required"`
	CreatedAt string `json:"created_at,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Content   string `json:"content,omitempty"`
}

type PastErrorResponse struct {
	Error string `json:"error"`
}

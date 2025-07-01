package model

import "github.com/gin-gonic/gin"

// omitempty : 비어있는 key 출력 x
type ResponseFormat struct {
	Message      string `json:"message,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Data         any    `json:"data,omitempty"`
}

func HandleResponse(c *gin.Context, statusCode int, code StatusCode, data any) {
	if statusCode >= 200 && statusCode < 300 {
		c.JSON(statusCode, ResponseFormat{
			Message: GetMessage(code),
			Data:    data,
		})
	} else {
		c.JSON(statusCode, ResponseFormat{
			ErrorMessage: GetMessage(code),
		})
	}
}

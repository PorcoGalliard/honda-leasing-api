package response

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorDetail{
			Message: message,
		},
	})
}

func ErrorWithCode(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
		},
	})
}


func NewValidationError(message string) error {
	return errors.New(message)
}

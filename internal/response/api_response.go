package response

import "github.com/gin-gonic/gin"

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendError(c *gin.Context, status int, msg string) {
	resp := ApiResponse{
		Success: false,
		Message: msg,
		Data:    nil,
	}
	c.JSON(status, resp)
}

package response

import "github.com/gin-gonic/gin"

type SuccessResponse struct {
	Data interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, SuccessResponse{Data: data})
}

func Error(c *gin.Context, status int, err string, details string) {
	c.JSON(status, ErrorResponse{
		Error:   err,
		Details: details,
	})
}

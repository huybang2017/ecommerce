package utils

import (
	resp "product-service/internal/dto/response"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, resp.ErrorResponse{
		Code:    code,
		Message: message,
	})
}

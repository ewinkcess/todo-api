package middleware

import (
	"net/http"
	"todo-api/internal/domain"
	"todo-api/internal/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last()
		if appErr, ok := err.Err.(utils.AppError); ok {
			c.JSON(appErr.Code, domain.NewErrorResponse(
				appErr.Message,
				appErr.Errors,
			))
			return
		}
		c.JSON(http.StatusInternalServerError, domain.NewErrorResponse(
			"Terjadi kesalahan pada server",
			err.Error(),
		))
	}
}

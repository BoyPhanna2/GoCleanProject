package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"myapp/internal/shared/domain"
	"myapp/internal/shared/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If there are errors attached to the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			var status int
			var errMsg = err.Error()
			var details = ""

			var ve validator.ValidationErrors
			switch {
			case errors.As(err, &ve):
				status = http.StatusBadRequest
				errMsg = "validation failed"
				var errorMessages []string
				for _, fe := range ve {
					errorMessages = append(errorMessages, getValidationMessage(fe))
				}
				details = strings.Join(errorMessages, ", ")
			case errors.Is(err, domain.ErrNotFound):
				status = http.StatusNotFound
			case errors.Is(err, domain.ErrConflict):
				status = http.StatusConflict
			case errors.Is(err, domain.ErrUnauthorized):
				status = http.StatusUnauthorized
			case errors.Is(err, domain.ErrBadRequest):
				status = http.StatusBadRequest
			default:
				status = http.StatusInternalServerError
			}

			response.Error(c, status, errMsg, details)
			// Stop further writing to response if it was already handled
			c.Abort()
		}
	}
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("'%s' is a required field", fe.Field())
	case "email":
		return fmt.Sprintf("'%s' must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("'%s' must be at least %s characters long", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("'%s' is invalid", fe.Field())
	}
}

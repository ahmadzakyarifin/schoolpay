package helper

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  bool                `json:"status"`
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors,omitempty"`
	Data    interface{}         `json:"data,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
	})
}

func ErrorResponseRaw(c *gin.Context, code int, err error) {
	c.JSON(code, Response{
		Status:  false,
		Message: err.Error(),
	})
}

func ErrorValidationResponse(c *gin.Context, code int, message string, errors map[string][]string) {
	c.JSON(code, Response{
		Status:  false,
		Message: message,
		Errors:  errors,
	})
}

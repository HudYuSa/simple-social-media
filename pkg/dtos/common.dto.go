package dtos

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// this package is where u put all data transfer object representation
// and all its function

type WebResponse struct {
	Message string `json:"message,omitempty"`
	Error   bool   `json:"error,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func RespondWithError(ctx *gin.Context, code int, errMsg string) {
	err := errors.New(errMsg)

	ctx.Error(err)
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, WebResponse{
		Error:   true,
		Message: errMsg,
	})
}

func RespondWithJson(ctx *gin.Context, code int, webResponse WebResponse) {
	ctx.JSON(code, webResponse)
}

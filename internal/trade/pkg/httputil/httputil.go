package httputil

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewError(ctx *gin.Context, status int, err string) {
	er := HTTPError{
		Code:    status,
		Message: err,
	}
	ctx.JSON(status, er)
}

func NewSuccess(ctx *gin.Context, data interface{}) {
	msg := "success"
	success := HTTPSuccess{
		Code:    http.StatusOK,
		Message: msg,
		Data:    data,
	}
	ctx.JSON(http.StatusOK, success)
}

type HTTPSuccess struct {
	Code    int         `json:"code" example:"200"`
	Message string      `json:"message" example:"success"`
	Data    interface{} `json:"data"`
}

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

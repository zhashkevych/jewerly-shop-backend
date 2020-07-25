package handler

import (
	"github.com/gin-gonic/gin"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

var (
	statusCodes = map[error]int{
		jewerly.ErrUserNotFound: http.StatusBadRequest,
	}
)

func getStatusCode(err error) int {
	code, ex := statusCodes[err]
	if !ex {
		return http.StatusInternalServerError
	}

	return code
}

func newErrorResponse(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, errorResponse{
		Error: err.Error(),
	})
}

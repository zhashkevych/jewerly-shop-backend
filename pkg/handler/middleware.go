package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AccessToken  = "Authorization"
)

func (h *Handler) adminIdentity(c *gin.Context) {
	header := c.Request.Header.Get(AccessToken)

	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("empty auth header"))
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid auth header"))
		return
	}

	if headerParts[1] == "" {
		newErrorResponse(c, http.StatusUnauthorized, errors.New("invalid token"))
		return
	}

	err := h.services.Admin.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}
}
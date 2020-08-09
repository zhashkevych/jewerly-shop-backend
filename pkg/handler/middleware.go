package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AccessToken  = "Authorization"
	UserCtx      = "user"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.Request.Header.Get(AccessToken)

	// Not Authorized
	if header == "" {
		c.Next()
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

	user, err := h.services.Auth.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.Set(UserCtx, user)
}

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
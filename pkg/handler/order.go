package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
)

func (h *Handler) placeOrder(c *gin.Context) {
	var inp jewerly.CreateOrderInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err := inp.Validate(); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	url, err := h.services.Order.Create(inp)
	if err != nil {
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"url": url,
	})
}

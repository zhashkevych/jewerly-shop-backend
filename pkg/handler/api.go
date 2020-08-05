package handler

import (
	"github.com/gin-gonic/gin"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
)

func (h *Handler) getUserProfile(c *gin.Context) {
	user, _ := c.Get(UserCtx)

	user, err := h.services.User.GetById(user.(jewerly.User).Id)
	if err != nil {
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getUserOrders(c *gin.Context) {
	user, _ := c.Get(UserCtx)

	user, err := h.services.User.GetById(user.(jewerly.User).Id)
	if err != nil {
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) placeOrder(c *gin.Context) {

}
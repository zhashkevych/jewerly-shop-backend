package handler

import (
	"github.com/gin-gonic/gin"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
)

func (h *Handler) getUserProfile(c *gin.Context) {
	user, _ := c.Get(UserCtx)

	c.JSON(http.StatusOK, user.(jewerly.User))
}

func (h *Handler) getUserOrders(c *gin.Context) {
	//user, _ := c.Get(UserCtx)
	//
	//user, err := h.services.User.GetById(user.(jewerly.User).Id)
	//if err != nil {
	//	newErrorResponse(c, getStatusCode(err), err)
	//	return
	//}
	//
	//c.JSON(http.StatusOK, user)
}

// NO AUTH NEEDED

func (h *Handler) placeOrder(c *gin.Context) {
	var inp jewerly.CreateOrderInput

	if err := c.ShouldBindJSON(&inp); err != nil {
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

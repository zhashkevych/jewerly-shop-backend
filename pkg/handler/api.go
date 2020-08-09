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

type placeOrderInput struct {
	ProductIds []int `json:"product_ids", binding:"required"`
}

func (h *Handler) placeOrder(c *gin.Context) {
	var inp placeOrderInput

	if err := c.ShouldBindJSON(&inp); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	user, ok := c.Get(UserCtx)
	if !ok {
		
	}

	err := h.services.Order.Create(user.(jewerly.User).Id, inp.ProductIds)
	if err != nil {
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.Status(http.StatusOK)
}
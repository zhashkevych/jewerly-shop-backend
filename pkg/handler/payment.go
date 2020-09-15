package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
)

func (h *Handler) callback(c *gin.Context) {
	var inp jewerly.TransactionCallbackInput

	if err := c.ShouldBind(&inp); err != nil {
		logrus.Errorf("failed to bind response: %s\n", err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	logrus.Debugf("input: %+v", inp)

	h.services.Order.ProcessCallback(inp)
	c.Status(http.StatusOK)
}

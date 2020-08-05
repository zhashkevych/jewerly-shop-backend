package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
)

// Products Handlers
func (h *Handler) createProduct(c *gin.Context) {
	var inp jewerly.CreateProductInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		logrus.Errorf("Failed to bind createProductInput structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := inp.Validate(); err != nil {
		logrus.Errorf("Failed to validate input body: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Product.Create(inp); err != nil {
		logrus.Errorf("Failed to create new product: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) updateProduct(c *gin.Context) {

}

func (h *Handler) deleteProduct(c *gin.Context) {

}

func (h *Handler) getAllProducts(c *gin.Context) {
	language := jewerly.GetLanguageFromQuery(c.Query("language"))

	products, err := h.services.Product.GetAll(jewerly.GetAllProductsFilters{
		Language: language,
	})
	if err != nil {
		logrus.Errorf("Failed to get products: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *Handler) getProduct(c *gin.Context) {

}

func (h *Handler) uploadImage(c *gin.Context) {

}

// Orders Handlers

func (h *Handler) getAllOrders(c *gin.Context) {

}

func (h *Handler) getOrder(c *gin.Context) {

}

package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
	"strconv"
)

const (
	maxUploadSize = 5 << 20 // 5 megabytes
)

var (
	imageTypes = map[string]interface{}{
		"image/jpeg": nil,
		"image/png":  nil,
	}
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("Failed to parse id from query: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Product.Delete(id); err != nil {
		logrus.Errorf("Failed to delete product: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.Status(http.StatusOK)
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
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("Failed to parse id from query: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	language := jewerly.GetLanguageFromQuery(c.Query("language"))

	product, err := h.services.Product.GetById(id, language)
	if err != nil {
		logrus.Errorf("Failed to delete product: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *Handler) uploadImage(c *gin.Context) {
	// Limit Upload File Size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

	file, fileHeader, err := c.Request.FormFile("image")
	if err != nil {
		logrus.Errorf("Failed to get image: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	defer file.Close()

	buffer := make([]byte, fileHeader.Size)
	file.Read(buffer)
	fileType := http.DetectContentType(buffer)

	// Validate File Type
	if _, ex := imageTypes[fileType]; !ex {
		logrus.Errorf("Failed to validate image type\n")
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid file type"))
		return
	}

	id, err := h.services.Product.UploadImage(c.Request.Context(), file, fileHeader.Size, fileType)
	if err != nil {
		logrus.Errorf("Failed to upload image: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// Orders Handlers

func (h *Handler) getAllOrders(c *gin.Context) {

}

func (h *Handler) getOrder(c *gin.Context) {

}

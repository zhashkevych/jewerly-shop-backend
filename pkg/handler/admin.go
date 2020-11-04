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

type adminSignInInput struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) adminSignIn(c *gin.Context) {
	var inp adminSignInInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		logrus.WithField("handler", "adminSignIn").Errorf("Failed to bind sign in structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	token, err := h.services.Admin.SignIn(inp.Login, inp.Password)
	if err != nil {
		logrus.WithField("handler", "adminSignIn").Errorf("Failed to sign in: %s\n", err.Error())
		newErrorResponse(c, http.StatusUnauthorized, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}

// Products Handlers
func (h *Handler) createProduct(c *gin.Context) {
	var inp jewerly.CreateProductInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		logrus.Errorf("Failed to bind createProductInput structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
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
	var inp jewerly.UpdateProductInput
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

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("Failed to parse id from query: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if err := h.services.Product.Update(id, inp); err != nil {
		logrus.Errorf("Failed to create new product: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) deleteProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("Failed to parse id from query: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	if id == 0 {
		logrus.Error("id is 0")
		newErrorResponse(c, http.StatusBadRequest, errors.New("id can't be zero"))
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
	products, err := h.services.Product.GetAll(getProductFilters(c))
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

	if id == 0 {
		logrus.Error("id is 0")
		newErrorResponse(c, http.StatusBadRequest, errors.New("id can't be zero"))
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
		newErrorResponse(c, http.StatusBadRequest, err)
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
	orders, err := h.services.Order.GetAll(getOrderFilters(c))
	if err != nil {
		logrus.Errorf("Failed to get orders: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *Handler) getOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Errorf("Failed to parse id from query: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	order, err := h.services.Order.GetById(id)
	if err != nil {
		logrus.Errorf("Failed to get order: %s\n", err.Error())
		newErrorResponse(c, getStatusCode(err), err)
		return
	}

	c.JSON(http.StatusOK, order)
}

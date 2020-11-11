package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"net/http"
	"strconv"
)

func (h *Handler) getHomepageImages(c *gin.Context) {
	images, err := h.services.Settings.GetImages()
	if err != nil {
		logrus.Errorf("Failed to get homepage images: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, images)
}

type homepageImageInput struct {
	ImageID int `json:"image_id" binding:"required"`
}

func (h *Handler) createHomepageImage(c *gin.Context) {
	var inp homepageImageInput
	if err := c.BindJSON(&inp); err != nil {
		logrus.Errorf("Failed to parse input body: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	err := h.services.Settings.CreateImage(inp.ImageID)
	if err != nil {
		logrus.Errorf("Failed to create homepage image: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) updateHomepageImage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Errorf("Failed to parse id param: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	var inp homepageImageInput
	if err := c.BindJSON(&inp); err != nil {
		logrus.Errorf("Failed to parse input body: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	err = h.services.Settings.UpdateImage(id, inp.ImageID)
	if err != nil {
		logrus.Errorf("Failed to update homepage image: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) getTextBlocks(c *gin.Context) {
	textBlocks, err := h.services.Settings.GetTextBlocks()
	if err != nil {
		logrus.Errorf("Failed to get text blocks: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, textBlocks)
}

func (h *Handler) getTextBlockById(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Errorf("Failed to parse id param: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	textBlock, err := h.services.Settings.GetTextBlockById(id)
	if err != nil {
		logrus.Errorf("Failed to get text blocks: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, textBlock)
}

func (h *Handler) createTextBlock(c *gin.Context) {
	var inp jewerly.TextBlock
	if err := c.BindJSON(&inp); err != nil {
		logrus.Errorf("Failed to parse input body: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	err := h.services.Settings.CreateTextBlock(inp)
	if err != nil {
		logrus.Errorf("Failed to create homepage image: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) updateTextBlock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		logrus.Errorf("Failed to parse id param: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid id param"))
		return
	}

	var inp jewerly.UpdateTextBlockInput
	if err := c.BindJSON(&inp); err != nil {
		logrus.Errorf("Failed to parse input body: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	if err := inp.Validate(); err != nil {
		logrus.Errorf("Validation afiled: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, errors.New("invalid input body"))
		return
	}

	err = h.services.Settings.UpdateTextBlock(id, inp)
	if err != nil {
		logrus.Errorf("Failed to update text block: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) getSettings(c *gin.Context) {
	settings, err := h.services.Settings.GetSettings()
	if err != nil {
		logrus.Errorf("Failed to get settings: %s\n", err.Error())
		newErrorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, settings)
}

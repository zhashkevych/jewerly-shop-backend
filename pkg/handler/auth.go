package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/service"
	"net/http"
)

type signUpInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var inp signUpInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		logrus.Errorf("Failed to bind signUp structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	err := h.services.Auth.SignUp(service.SignUpInput{
		FirstName: inp.FirstName,
		LastName:  inp.LastName,
		Email:     inp.Email,
		Password:  inp.Password,
	})
	if err != nil {
		logrus.Errorf("Failed to create user: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.Status(http.StatusOK)
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var inp signInInput
	if err := c.ShouldBindJSON(&inp); err != nil {
		logrus.Errorf("Failed to bind signUp structure: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	token, err := h.services.Auth.SignIn(inp.Email, inp.Password)
	if err != nil {
		logrus.Errorf("Failed to create user: %s\n", err.Error())
		newErrorResponse(c, http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, signInResponse{
		Token: token,
	})
}

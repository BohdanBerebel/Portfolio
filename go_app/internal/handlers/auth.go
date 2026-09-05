package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/BohdanBerebel/Portfolio/go_app/internal/services"
	"github.com/BohdanBerebel/Portfolio/go_app/pkg/response"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	userService *services.UserService
}

func NewAuthHandler(userService *services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	result, err := h.userService.Register(
		c.Request.Context(),
		input.Email,
		input.Password,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.JSON(
		c,
		http.StatusCreated,
		gin.H{
			"id":    result.User.ID,
			"email": result.User.Email,
			"token": result.Token,
		},
	)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	result, err := h.userService.Login(
		c.Request.Context(),
		input.Email,
		input.Password,
	)
	if err != nil {
		handleError(c, err)
		return
	}

	response.JSON(
		c,
		http.StatusOK,
		gin.H{
			"id":    result.User.ID,
			"email": result.User.Email,
			"token": result.Token,
		},
	)

}

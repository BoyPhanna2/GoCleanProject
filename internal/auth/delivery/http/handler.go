package http

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"myapp/internal/auth/domain"
	"myapp/internal/shared/response"
)

type AuthHandler struct {
	usecase domain.AuthUsecase
}

func NewAuthHandler(r *gin.Engine, usecase domain.AuthUsecase) {
	handler := &AuthHandler{usecase: usecase}

	authGroup := r.Group("/api/v1/auth")
	{
		authGroup.POST("/register", handler.Register)
		authGroup.POST("/login", handler.Login)
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err) // Centralized error handling will pick this up
		return
	}

	user, err := h.usecase.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
	})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	token, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, gin.H{
		"token": token,
	})
}

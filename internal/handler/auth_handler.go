package handler

import (
	"net/http"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
	"github.com/YannisMuminov/internal/middleware"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service domain.AuthService
}

func NewAuthHandler(service domain.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, tokens)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req domain.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refresh, err := h.service.Refresh(c.Request.Context(), &req)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, refresh)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.service.Me(c.Request.Context(), userID)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	perms := make([]string, len(user.Permissions))

	for i, p := range user.Permissions {
		perms[i] = p.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          user.ID,
		"email":       user.Email,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"is_active":   user.IsActive,
		"role":        user.Role.Name,
		"permissions": perms,
	})
}

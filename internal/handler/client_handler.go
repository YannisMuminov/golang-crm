package handler

import (
	"net/http"
	"strconv"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
	"github.com/YannisMuminov/internal/middleware"
	"github.com/gin-gonic/gin"
)

type ClientHandler struct {
	service domain.ClientService
}

func NewClientHandler(service domain.ClientService) *ClientHandler {
	return &ClientHandler{service: service}
}

func (h *ClientHandler) Create(c *gin.Context) {
	var req domain.ClientCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	client, err := h.service.Create(c.Request.Context(), &req, userID)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func (h *ClientHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) GetAll(c *gin.Context) {
	filter := domain.ClientFilter{
		Search: c.Query("search"),
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 20),
	}

	client, total, err := h.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": client,
		"meta": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (h *ClientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req domain.ClientUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "client deleted"})
}

func queryInt(c *gin.Context, key string, defaultValue int) int {
	val := c.Query(key)
	if val == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}
	return i
}

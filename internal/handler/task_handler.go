package handler

import (
	"net/http"
	"strconv"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
	"github.com/YannisMuminov/internal/middleware"
	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service domain.TaskService
}

func NewTaskHandler(service domain.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req domain.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	task, err := h.service.Create(c.Request.Context(), &req, userID)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "invalid id"})
		return
	}

	task, err := h.service.GetByID(c.Request.Context(), id)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetAll(c *gin.Context) {
	filter := domain.TaskFilter{
		DealID:     queryInt64(c, "deal_id"),
		AssignedTo: queryInt64(c, "assigned_id"),
		Status:     domain.TaskStatus(c.Query("status")),
		Page:       queryInt(c, "page", 1),
		Limit:      queryInt(c, "limit", 20),
	}

	tasks, total, err := h.service.GetAll(c.Request.Context(), filter)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"meta": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.service.Update(c.Request.Context(), id, &req)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

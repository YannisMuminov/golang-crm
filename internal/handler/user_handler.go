package handler

import (
	"net/http"
	"strconv"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (u *UserHandler) GetAll(c *gin.Context) {
	filter := domain.UserFilter{
		RoleID: queryInt64(c, "role_id"),
		Page:   queryInt(c, "page", 1),
		Limit:  queryInt(c, "limit", 20),
	}

	if val := c.Query("is_active"); val != "" {
		b, err := strconv.ParseBool(val)
		if err != nil {
			filter.IsActive = &b
		}
	}

	users, total, err := u.service.GetAll(c.Request.Context(), filter)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"meta": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (u *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.GetByID(c.Request.Context(), id)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *UserHandler) Deactivate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := u.service.Deactivate(c.Request.Context(), id); err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deactivated"})
}

func (u *UserHandler) Activate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := u.service.Activate(c.Request.Context(), id); err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user activated"})
}

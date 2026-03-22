package handler

import (
	"net/http"
	"strconv"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
	"github.com/YannisMuminov/internal/middleware"
	"github.com/gin-gonic/gin"
)

type DealHandler struct {
	service domain.DealService
}

func NewDealHandler(service domain.DealService) *DealHandler {
	return &DealHandler{service: service}
}

func (d *DealHandler) Create(c *gin.Context) {
	var req domain.CreateDealRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := middleware.GetUserID(c)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deal, err := d.service.Create(c.Request.Context(), &req, userID)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, deal)
}

func (d *DealHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	deal, err := d.service.GetByID(c.Request.Context(), id)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, deal)
}

func (d *DealHandler) GetAll(c *gin.Context) {
	filter := domain.DealFilter{
		Status:     domain.DealStatus(c.Query("status")),
		ClientID:   queryInt64(c, "client_id"),
		AssignedTo: queryInt64(c, "assigned_to"),
		Page:       queryInt(c, "page", 1),
		Limit:      queryInt(c, "limit", 20),
	}

	deals, total, err := d.service.GetAll(c.Request.Context(), filter)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": deals,
		"meta": gin.H{
			"total": total,
			"page":  filter.Page,
			"limit": filter.Limit,
		},
	})
}

func (d *DealHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req domain.UpdateDealRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deal, err := d.service.Update(c.Request.Context(), id, &req)

	if err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, deal)
}

func (d *DealHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := d.service.Delete(c.Request.Context(), id); err != nil {
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deal deleted"})
}

func queryInt64(c *gin.Context, key string) int64 {
	val := c.Query(key)
	if val == "" {
		return 0
	}
	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

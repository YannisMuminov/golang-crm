package handler

import (
	"fmt"
	"net/http"

	"github.com/YannisMuminov/internal/apperror"
	"github.com/YannisMuminov/internal/domain"
	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	service domain.StatsService
}

func NewStatsHandler(service domain.StatsService) *StatsHandler {
	return &StatsHandler{service: service}
}

func (s *StatsHandler) GetStats(c *gin.Context) {
	stats, err := s.service.GetStats(c.Request.Context())

	if err != nil {
		fmt.Println("httpErrrrrr", err)
		code, msg := apperror.HTTPStatus(err)
		c.JSON(code, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, stats)
}

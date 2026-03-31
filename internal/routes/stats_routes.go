package routes

import (
	"github.com/YannisMuminov/internal/config"
	"github.com/YannisMuminov/internal/handler"
	"github.com/YannisMuminov/internal/middleware"
	"github.com/YannisMuminov/internal/service"
	"github.com/YannisMuminov/pkg/database"
	"github.com/gin-gonic/gin"
)

func StatsRoutes(r *gin.Engine, cfg *config.Config) {

	statsService := service.NewStatsService(database.GetDB())
	statsHandler := handler.NewStatsHandler(statsService)

	stats := r.Group("/stats")

	stats.Use(middleware.RequireAuth(&cfg.JWT))
	{
		stats.GET("", middleware.RequirePermission("user:read"), statsHandler.GetStats)
	}

}

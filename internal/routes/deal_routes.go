package routes

import (
	"github.com/YannisMuminov/internal/config"
	"github.com/YannisMuminov/internal/handler"
	"github.com/YannisMuminov/internal/middleware"
	"github.com/YannisMuminov/internal/repository"
	"github.com/YannisMuminov/internal/service"
	"github.com/YannisMuminov/pkg/database"
	"github.com/gin-gonic/gin"
)

func DealRoutes(r *gin.Engine, cfg *config.Config) {
	clientRepo := repository.NewClientRepository(database.GetDB())
	dealRepo := repository.NewDealRepository(database.GetDB())
	dealService := service.NewDealService(dealRepo, clientRepo)
	dealHandler := handler.NewDealHandler(dealService)

	deals := r.Group("/deals")
	deals.Use(middleware.RequireAuth(&cfg.JWT))
	{
		deals.GET("", middleware.RequirePermission("deals:read"), dealHandler.GetAll)
		deals.GET("/:id", middleware.RequirePermission("deals:read"), dealHandler.GetByID)
		deals.POST("", middleware.RequirePermission("deals:write"), dealHandler.Create)
		deals.PUT("/:id", middleware.RequirePermission("deals:write"), dealHandler.Update)
		deals.DELETE("/:id", middleware.RequirePermission("deals:delete"), dealHandler.Delete)
	}

}

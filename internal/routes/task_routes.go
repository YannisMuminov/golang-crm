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

func TaskRoutes(r *gin.Engine, cfg *config.Config) {

	dealRepo := repository.NewDealRepository(database.GetDB())
	taskRepo := repository.NewTaskRepository(database.GetDB())
	taskSvc := service.NewTaskService(taskRepo, dealRepo)
	taskHandler := handler.NewTaskHandler(taskSvc)

	tasks := r.Group("/tasks")
	tasks.Use(middleware.RequireAuth(&cfg.JWT))
	{
		tasks.GET("", middleware.RequirePermission("tasks:read"), taskHandler.GetAll)
		tasks.GET("/:id", middleware.RequirePermission("tasks:read"), taskHandler.GetByID)
		tasks.POST("", middleware.RequirePermission("tasks:write"), taskHandler.Create)
		tasks.PUT("/:id", middleware.RequirePermission("tasks:write"), taskHandler.Update)
		tasks.DELETE("/:id", middleware.RequirePermission("tasks:delete"), taskHandler.Delete)
	}
}

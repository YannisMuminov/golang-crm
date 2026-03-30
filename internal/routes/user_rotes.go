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

func UserRoutes(r *gin.Engine, cfg *config.Config) {
	userRepo := repository.NewUserRepository(database.GetDB())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	users := r.Group("/users")

	users.Use(middleware.RequireAuth(&cfg.JWT))
	{
		users.GET("", middleware.RequirePermission("users:read"), userHandler.GetAll)
		users.GET("/:id", middleware.RequirePermission("users:read"), userHandler.GetByID)
		users.PUT("/:id", middleware.RequirePermission("users:write"), userHandler.Update)
		users.PUT("/:id/deactivate", middleware.RequirePermission("users:write"), userHandler.Deactivate)
		users.PUT("/:id/activate", middleware.RequirePermission("users:write"), userHandler.Activate)
	}
}

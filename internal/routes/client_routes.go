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

func ClientRoutes(r *gin.Engine, cfg *config.Config) {
	clientRepo := repository.NewClientRepository(database.GetDB())
	clientSrv := service.NewClientService(clientRepo)
	clientHandler := handler.NewClientHandler(clientSrv)

	clients := r.Group("/clients")
	clients.Use(middleware.RequireAuth(&cfg.JWT))
	{
		clients.GET("", middleware.RequirePermission("clients:read"), clientHandler.GetAll)
		clients.GET("/:id", middleware.RequirePermission("clients:read"), clientHandler.GetByID)
		clients.POST("", middleware.RequirePermission("clients:write"), clientHandler.Create)
		clients.PUT("/:id", middleware.RequirePermission("clients:write"), clientHandler.Update)
		clients.DELETE("/:id", middleware.RequirePermission("clients:delete"), clientHandler.Delete)
	}
}

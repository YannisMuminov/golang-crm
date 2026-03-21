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

func AuthRoutes(r *gin.Engine, cfg *config.Config) {

	authRepo := repository.NewAuthRepository(database.GetDB())
	authService := service.NewAuthService(authRepo, &cfg.JWT)
	authHandler := handler.NewAuthHandler(authService)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	protected := r.Group("")
	protected.Use(middleware.RequireAuth(&cfg.JWT))
	{
		protected.GET("/me", authHandler.Me)
	}
}

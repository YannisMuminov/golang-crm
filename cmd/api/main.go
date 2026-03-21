package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YannisMuminov/internal/config"
	"github.com/YannisMuminov/internal/routes"
	"github.com/YannisMuminov/pkg/database"
	"github.com/YannisMuminov/pkg/database/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.LoadConfig()

	if err := logger.Init(cfg.Server.AppEnv); err != nil {
		panic("failed to init logger: " + err.Error())
	}
	defer logger.Sync()

	if err := database.InitDB(&cfg.Database); err != nil {
		logger.L.Fatal("failed to connected database", zap.Error(err))
	}
	defer database.CloseDB()

	var r *gin.Engine

	if cfg.Server.IsDevelopment() {
		r = gin.New()
		r.Use(gin.Recovery())
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
		r.Use(gin.Recovery())
	}

	r.Use(logger.GinLogger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
		MaxAge:        12 * time.Hour,
	}))

	routes.AuthRoutes(r, cfg)
	routes.ClientRoutes(r, cfg)

	addr := fmt.Sprintf(":%s", cfg.Server.AppPort)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.L.Info("Server starting",
			zap.String("addr", "http://localhost"+addr),
			zap.String("env", cfg.Server.AppEnv),
		)
		if err := srv.ListenAndServe(); err != nil || errors.Is(err, http.ErrServerClosed) {
			logger.L.Fatal("server failed", zap.Error(err))
		}

	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.L.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.L.Fatal("forced shutdown", zap.Error(err))
	}
	logger.L.Info("✅ server exited cleanly")
}

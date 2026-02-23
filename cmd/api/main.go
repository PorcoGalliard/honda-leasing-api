package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nanasuryana335/honda-leasing-api/cmd/api/routes"
	configs "github.com/nanasuryana335/honda-leasing-api/internal/config"
	postgres "github.com/nanasuryana335/honda-leasing-api/pkg"
)

func main() {
	config := configs.Load()
	db, err := postgres.InitDB(config)
	if err != nil {
		log.Fatalf("Failed to initialize database")
	}

	defer postgres.CloseDB(db)

	// Set Gin mode based on environment
	if config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Setup routes
	router := gin.Default()
	routes.SetupRoutes(router, db.DB, config)

	addr := config.Server.Address
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	// Start server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server starting on %s in %s mode", addr, config.Environment)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exiting")
}

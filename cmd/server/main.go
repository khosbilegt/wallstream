package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.io/khosbilegt/wallstream/internal/server/api"
	"github.io/khosbilegt/wallstream/internal/server/api/handlers"
	"github.io/khosbilegt/wallstream/internal/server/db"
	"github.io/khosbilegt/wallstream/internal/server/repository"
	"github.io/khosbilegt/wallstream/internal/server/service"
)

func main() {
	// Get MongoDB URI from environment or use default
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017/wallstream"
	}

	// Get server port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Load templates
	api.LoadTemplates()

	// Connect to MongoDB
	log.Println("Connecting to MongoDB...")
	client, err := db.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	// Get database
	database := client.Database("wallpaper-share")

	// Initialize collections
	collections := db.NewCollections(database)

	// Initialize repositories
	usersRepo := repository.NewUsersRepository(collections.Users)
	publisherRepo := repository.NewPublisherDeviceRepository(collections.PublisherDevices)
	publishedWallpaperRepo := repository.NewPublishedWallpaperRepository(collections.PublishedWallpapers)

	// Initialize services
	usersService := service.NewUsersService(usersRepo)

	fileService := service.NewFileService("uploads")
	publisherService := service.NewPublisherService(publisherRepo, publishedWallpaperRepo)

	// Initialize handlers
	handlers := handlers.NewHandlers(handlers.NewUserHandlers(usersService), handlers.NewFileHandlers(fileService), handlers.NewPublisherHandlers(publisherService))

	// Setup routes with Chi router
	router := chi.NewRouter()
	routes := api.NewRoutes(router, handlers)
	routes.RegisterRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

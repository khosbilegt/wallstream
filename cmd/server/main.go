package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.io/khosbilegt/wallstream/internal/server"
)

func main() {
	var (
		addr          = flag.String("addr", ":8080", "Server address")
		wallpapersDir = flag.String("wallpapers", "./wallpapers", "Directory to store wallpapers")
		baseURL       = flag.String("base-url", "http://localhost:8080", "Base URL for wallpaper links")
	)
	flag.Parse()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	srv, err := server.NewServer(*wallpapersDir, *baseURL, rdb)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting server on %s", *addr)
	log.Printf("Wallpapers directory: %s", *wallpapersDir)
	log.Printf("Base URL: %s", *baseURL)

	httpServer := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down server...")
	// In a production server, you'd use httpServer.Shutdown(context.Background())
	log.Println("Server stopped")
}

package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.io/khosbilegt/wallstream/internal/server/api/handlers"
)

type Routes struct {
	r        chi.Router
	handlers *handlers.Handlers
}

func NewRoutes(r chi.Router, handlers *handlers.Handlers) *Routes {
	return &Routes{r: r, handlers: handlers}
}

func (rts *Routes) RegisterRoutes() {
	// Apply global middleware
	rts.r.Use(middleware.RequestID)
	rts.r.Use(middleware.RealIP)
	rts.r.Use(middleware.Logger)
	rts.r.Use(middleware.Recoverer)

	// Public routes
	rts.r.Group(func(r chi.Router) {
		rts.r.Post("/api/users/register", rts.handlers.UserHandlers.CreateUser)
	})

	// File routes
	rts.r.Group(func(r chi.Router) {
		r.Use(rts.handlers.AuthMiddleware)
		r.Post("/api/files/upload", rts.handlers.FileHandlers.UploadWallpaper)
	})

	// Protected routes (API key authentication)
	rts.r.Group(func(r chi.Router) {
		r.Use(rts.handlers.AuthMiddleware)
		r.Post("/api/publisher/devices", rts.handlers.PublisherHandlers.CreatePublisherDevice)
		r.Get("/api/publisher/devices", rts.handlers.PublisherHandlers.GetPublisherDevices)
		r.Get("/api/publisher/devices/{deviceID}", rts.handlers.PublisherHandlers.GetPublisherDeviceByDeviceID)
		r.Delete("/api/publisher/devices/{deviceID}", rts.handlers.PublisherHandlers.DeletePublisherDeviceByDeviceID)
		r.Get("/api/publisher/devices/{deviceID}/upload-url", rts.handlers.PublisherHandlers.GetUploadURL)
		r.Post("/api/publisher/wallpaper", rts.handlers.PublisherHandlers.PublishUploadedWallpaper)
		r.Get("/api/publisher/wallpaper", rts.handlers.PublisherHandlers.GetPublishedWallpapers)
		r.Get("/api/publisher/wallpaper/{deviceID}", rts.handlers.PublisherHandlers.GetPublishedWallpapersByDeviceID)
		r.Delete("/api/publisher/wallpaper/{hash}", rts.handlers.PublisherHandlers.DeletePublishedWallpaperByHash)
		r.Get("/api/wallpaper/{deviceID}", rts.handlers.PublisherHandlers.ServeWallpaper)
	})
}

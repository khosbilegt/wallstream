package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Routes struct {
	r        chi.Router
	handlers *Handlers
}

func NewRoutes(r chi.Router, handlers *Handlers) *Routes {
	return &Routes{r: r, handlers: handlers}
}

func (rts *Routes) RegisterRoutes() {
	// Apply global middleware
	rts.r.Use(middleware.RequestID)
	rts.r.Use(middleware.RealIP)
	rts.r.Use(middleware.Logger)
	rts.r.Use(middleware.Recoverer)

	// Web routes
	rts.r.Get("/", rts.handlers.WebIndex)

	// Public routes
	rts.r.Post("/api/users/register", rts.handlers.CreateUser)

	// Protected routes (API key authentication)
	rts.r.Group(func(r chi.Router) {
		r.Use(rts.handlers.AuthMiddleware)
		r.Post("/api/files/upload", rts.handlers.UploadWallpaper)

		r.Post("/api/publisher/devices", rts.handlers.CreatePublisherDevice)
		r.Get("/api/publisher/devices", rts.handlers.GetPublisherDevices)
		r.Get("/api/publisher/devices/{deviceID}", rts.handlers.GetPublisherDeviceByDeviceID)
		r.Delete("/api/publisher/devices/{deviceID}", rts.handlers.DeletePublisherDeviceByDeviceID)
		r.Get("/api/publisher/devices/{deviceID}/upload-url", rts.handlers.GetUploadURL)
		r.Post("/api/publisher/wallpaper", rts.handlers.PublishUploadedWallpaper)
		r.Get("/api/publisher/wallpaper", rts.handlers.GetPublishedWallpapers)
		r.Get("/api/publisher/wallpaper/{deviceID}", rts.handlers.GetPublishedWallpapersByDeviceID)
	})
}

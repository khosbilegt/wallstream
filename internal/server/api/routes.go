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

		// Publisher
		r.Get("/api/publisher/state", rts.handlers.GetPublisherState)
		r.Post("/api/publisher/state", rts.handlers.CreatePublisherState)

		// Subscriber
		r.Get("/api/subscriber/state", rts.handlers.GetSubscriberState)
		r.Post("/api/subscriber/state", rts.handlers.CreateSubscriberState)
	})
}

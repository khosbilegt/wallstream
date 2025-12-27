package api

import "net/http"

type Routes struct {
	mux      *http.ServeMux
	handlers *Handlers
}

func NewRoutes(mux *http.ServeMux, handlers *Handlers) *Routes {
	return &Routes{mux: mux, handlers: handlers}
}

func (r *Routes) RegisterRoutes() {
	// Public routes
	r.mux.HandleFunc("/users/register", r.handlers.RegisterUser)
	r.mux.HandleFunc("/users/login", r.handlers.LoginUser)
	r.mux.HandleFunc("/users/refresh", r.handlers.RefreshToken)

	// Protected routes (require authentication)
	r.mux.HandleFunc("/publisher/state", r.handlers.AuthMiddleware(r.handlers.GetPublisherState))
	r.mux.HandleFunc("/subscriber/state", r.handlers.AuthMiddleware(r.handlers.GetSubscriberState))
	r.mux.HandleFunc("/publisher/state", r.handlers.AuthMiddleware(r.handlers.CreatePublisherState))
	r.mux.HandleFunc("/subscriber/state", r.handlers.AuthMiddleware(r.handlers.CreateSubscriberState))
}

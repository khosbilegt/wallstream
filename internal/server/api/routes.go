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
	r.mux.HandleFunc("/publisher/state", r.handlers.GetPublisherState)
	r.mux.HandleFunc("/subscriber/state", r.handlers.GetSubscriberState)
	r.mux.HandleFunc("/publisher/state", r.handlers.CreatePublisherState)
	r.mux.HandleFunc("/subscriber/state", r.handlers.CreateSubscriberState)
}

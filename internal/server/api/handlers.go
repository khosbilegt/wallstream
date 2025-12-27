package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.io/khosbilegt/wallstream/internal/server/repository"
	"github.io/khosbilegt/wallstream/internal/server/service"
)

type Handlers struct {
	publisherService  *service.PublisherService
	subscriberService *service.SubscriberService
}

func (h *Handlers) GetPublisherState(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.URL.Query().Get("id")
	state, err := h.publisherService.GetPublisherState(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(state)
}

func (h *Handlers) GetSubscriberState(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.URL.Query().Get("id")
	state, err := h.subscriberService.GetSubscriberState(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(state)
}

func (h *Handlers) CreatePublisherState(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var state repository.PublisherState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.publisherService.CreatePublisherState(ctx, &state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) CreateSubscriberState(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var state repository.SubscriberState
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = h.subscriberService.CreateSubscriberState(ctx, &state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

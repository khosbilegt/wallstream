package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.io/khosbilegt/wallstream/internal/server/repository"
	"github.io/khosbilegt/wallstream/internal/server/service"
)

type Handlers struct {
	usersService      *service.UsersService
	publisherService  *service.PublisherService
	subscriberService *service.SubscriberService
}

func NewHandlers(usersService *service.UsersService, publisherService *service.PublisherService, subscriberService *service.SubscriberService) *Handlers {
	return &Handlers{
		usersService:      usersService,
		publisherService:  publisherService,
		subscriberService: subscriberService,
	}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Define a struct matching the expected JSON
	var req struct {
		Username string `json:"username"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	apiKey, err := h.usersService.CreateUser(context.Background(), req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"username": req.Username, "api_key": apiKey})
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

func (h *Handlers) WebIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Wallcast Dashboard",
	}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

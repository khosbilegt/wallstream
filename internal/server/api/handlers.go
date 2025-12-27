package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

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

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists
	existingUser, _ := h.usersService.GetUserByUsername(ctx, user.Username)
	if existingUser != nil {
		http.Error(w, "username already exists", http.StatusConflict)
		return
	}

	// Set timestamps
	now := time.Now().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now

	err = h.usersService.CreateUser(ctx, &user, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate API key
	apiKey, err := service.GenerateAPIKey()
	if err != nil {
		http.Error(w, "failed to generate API key", http.StatusInternalServerError)
		return
	}

	// Save API key to user
	err = h.usersService.UpdateUserAPIKey(ctx, user.ID, apiKey)
	if err != nil {
		http.Error(w, "failed to save API key", http.StatusInternalServerError)
		return
	}

	// Return user with API key (don't send password)
	user.Password = ""
	user.APIKey = apiKey
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var loginReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.usersService.LoginUser(ctx, loginReq.Username, loginReq.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate new API key if user doesn't have one
	if user.APIKey == "" {
		apiKey, err := service.GenerateAPIKey()
		if err != nil {
			http.Error(w, "failed to generate API key", http.StatusInternalServerError)
			return
		}
		err = h.usersService.UpdateUserAPIKey(ctx, user.ID, apiKey)
		if err != nil {
			http.Error(w, "failed to save API key", http.StatusInternalServerError)
			return
		}
		user.APIKey = apiKey
	}

	// Don't send password back
	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

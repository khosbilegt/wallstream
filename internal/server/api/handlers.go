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
	jwtService        *service.JWTService
	refreshTokensRepo *repository.RefreshTokensRepository
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

	err = h.usersService.CreateUser(ctx, &user, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Save refresh token to database
	refreshToken := &repository.RefreshToken{
		Token:     tokenPair.RefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		CreatedAt: time.Now().Unix(),
	}
	if err := h.refreshTokensRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		http.Error(w, "failed to save refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
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

	// Generate JWT tokens
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Delete old refresh tokens for this user (optional: for single device login)
	// h.refreshTokensRepo.DeleteUserRefreshTokens(ctx, user.ID)

	// Save refresh token to database
	refreshToken := &repository.RefreshToken{
		Token:     tokenPair.RefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		CreatedAt: time.Now().Unix(),
	}
	if err := h.refreshTokensRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		http.Error(w, "failed to save refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

func (h *Handlers) RefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var refreshReq struct {
		RefreshToken string `json:"refresh_token"`
	}
	err := json.NewDecoder(r.Body).Decode(&refreshReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate refresh token
	claims, err := h.jwtService.ValidateRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Check if refresh token exists in database
	storedToken, err := h.refreshTokensRepo.GetRefreshToken(ctx, refreshReq.RefreshToken)
	if err != nil {
		http.Error(w, "refresh token not found", http.StatusUnauthorized)
		return
	}

	// Check if token is expired
	if time.Now().Unix() > storedToken.ExpiresAt {
		// Clean up expired token
		h.refreshTokensRepo.DeleteRefreshToken(ctx, refreshReq.RefreshToken)
		http.Error(w, "refresh token expired", http.StatusUnauthorized)
		return
	}

	// Get user to verify they still exist
	user, err := h.usersService.GetUserByID(ctx, claims.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	}

	// Generate new token pair
	tokenPair, err := h.jwtService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		http.Error(w, "failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Delete old refresh token
	h.refreshTokensRepo.DeleteRefreshToken(ctx, refreshReq.RefreshToken)

	// Save new refresh token
	refreshToken := &repository.RefreshToken{
		Token:     tokenPair.RefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		CreatedAt: time.Now().Unix(),
	}
	if err := h.refreshTokensRepo.SaveRefreshToken(ctx, refreshToken); err != nil {
		http.Error(w, "failed to save refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenPair)
}

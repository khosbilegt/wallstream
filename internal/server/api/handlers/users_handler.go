package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.io/khosbilegt/wallstream/internal/server/service"
	"github.io/khosbilegt/wallstream/internal/server/utils"
)

type UserHandlers struct {
	usersService *service.UsersService
}

func NewUserHandlers(usersService *service.UsersService) *UserHandlers {
	return &UserHandlers{usersService: usersService}
}

func (h *UserHandlers) CreateUser(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	apiKey, err := h.usersService.CreateUser(context.Background(), req.Username)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"username": req.Username, "api_key": apiKey})
}

// func (h *UserHandlers) WebIndex(w http.ResponseWriter, r *http.Request) {
// 	data := struct {
// 		Title string
// 	}{
// 		Title: "Wallcast Dashboard",
// 	}
// 	err := templates.ExecuteTemplate(w, "index.html", data)
// 	if err != nil {
// 		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
// 			"error": err.Error(),
// 		})
// 	}
// }

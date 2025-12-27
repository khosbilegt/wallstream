package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.io/khosbilegt/wallstream/internal/server/service"
)

type Handlers struct {
	usersService     *service.UsersService
	fileService      *service.FileService
	publisherService *service.PublisherService
}

func NewHandlers(usersService *service.UsersService, fileService *service.FileService, publisherService *service.PublisherService) *Handlers {
	return &Handlers{
		usersService:     usersService,
		fileService:      fileService,
		publisherService: publisherService,
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

func (h *Handlers) CreatePublisherDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	// Define a struct matching the expected JSON
	var req struct {
		DeviceID string `json:"device_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	// Check if the device ID is already in use
	publisherDevice, err := h.publisherService.GetPublisherDeviceByDeviceID(r.Context(), req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if publisherDevice != nil {
		http.Error(w, "device ID already in use", http.StatusBadRequest)
		return
	}

	err = h.publisherService.CreatePublisherDevice(r.Context(), userID, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"device_id": req.DeviceID})
}

func (h *Handlers) GetPublisherDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	publisherDevices, err := h.publisherService.GetPublisherDevicesByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(publisherDevices)
}

func (h *Handlers) GetPublisherDeviceByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	publisherDevice, err := h.publisherService.GetPublisherDeviceByDeviceID(r.Context(), deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(publisherDevice)
}

func (h *Handlers) DeletePublisherDeviceByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	err := h.publisherService.DeletePublisherDeviceByDeviceID(r.Context(), deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Upload wallpaper to the server
func (h *Handlers) UploadWallpaper(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename, err := h.fileService.UploadFileStream(r.Context(), file, fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"filename": filename})

}

func (h *Handlers) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	uploadURL, err := h.publisherService.GenerateUploadURL(r.Context(), userID, deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"upload_url": uploadURL})
}

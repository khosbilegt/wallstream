package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.io/khosbilegt/wallstream/internal/server/service"
	"github.io/khosbilegt/wallstream/internal/server/utils"
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
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	apiKey, err := h.usersService.CreateUser(context.Background(), req.Username)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"username": req.Username, "api_key": apiKey})
}

func (h *Handlers) WebIndex(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
	}{
		Title: "Wallcast Dashboard",
	}
	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
}

func (h *Handlers) CreatePublisherDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
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
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	// Check if the device ID is already in use
	publisherDevice, err := h.publisherService.GetPublisherDeviceByDeviceID(r.Context(), req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if publisherDevice != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "device ID already in use",
		})
		return
	}

	err = h.publisherService.CreatePublisherDevice(r.Context(), userID, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"device_id": req.DeviceID})
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
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, publisherDevices)
}

func (h *Handlers) GetPublisherDeviceByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	publisherDevice, err := h.publisherService.GetPublisherDeviceByDeviceID(r.Context(), deviceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, publisherDevice)
}

func (h *Handlers) DeletePublisherDeviceByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	err := h.publisherService.DeletePublisherDeviceByDeviceID(r.Context(), deviceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
}

func (h *Handlers) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
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

	writeJSON(w, http.StatusOK, map[string]string{"upload_url": uploadURL})
}

// Upload wallpaper to the server
func (h *Handlers) UploadWallpaper(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	filename, err := h.fileService.UploadFileStream(r.Context(), file, fileHeader.Filename)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"filename": filename})

}
func (h *Handlers) PublishUploadedWallpaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Validate user ID
	userID, ok := utils.GetStringFromContext(r.Context(), "user_id")
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing or invalid user_id",
		})
		return
	}

	var req struct {
		Filename string `json:"filename"`
		DeviceID string `json:"device_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	if err := h.publisherService.PublishUploadedWallpaper(
		r.Context(),
		userID,
		req.DeviceID,
		req.Filename,
	); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Wallpaper published successfully",
	})
}

func (h *Handlers) GetPublishedWallpapers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	publishedWallpapers, err := h.publisherService.GetPublishedWallpapersByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, publishedWallpapers)
}

func (h *Handlers) GetPublishedWallpapersByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	userID, ok := utils.GetStringFromContext(r.Context(), "user_id")
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing or invalid user_id",
		})
		return
	}

	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing deviceID",
		})
		return
	}

	publishedWallpapers, err := h.publisherService.
		GetPublishedWallpapersByDeviceID(r.Context(), userID, deviceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, publishedWallpapers)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

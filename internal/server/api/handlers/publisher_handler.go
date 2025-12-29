package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.io/khosbilegt/wallstream/internal/server/service"
	"github.io/khosbilegt/wallstream/internal/server/utils"
)

type PublisherHandlers struct {
	publisherService *service.PublisherService
}

func NewPublisherHandlers(publisherService *service.PublisherService) *PublisherHandlers {
	return &PublisherHandlers{publisherService: publisherService}
}

func (h *PublisherHandlers) CreatePublisherDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
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
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
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
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "device ID already in use",
		})
		return
	}

	err = h.publisherService.CreatePublisherDevice(r.Context(), userID, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"device_id": req.DeviceID})
}

func (h *PublisherHandlers) GetPublisherDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	publisherDevices, err := h.publisherService.GetPublisherDevicesByUserID(r.Context(), userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, publisherDevices)
}

func (h *PublisherHandlers) GetPublisherDeviceByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	publisherDevice, err := h.publisherService.GetPublisherDeviceByDeviceID(r.Context(), deviceID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, publisherDevice)
}

func (h *PublisherHandlers) DeletePublisherDeviceByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the device ID from the request context
	deviceID := r.Context().Value("device_id").(string)

	err := h.publisherService.DeletePublisherDeviceByDeviceID(r.Context(), deviceID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
}

func (h *PublisherHandlers) GetUploadURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
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

	utils.WriteJSON(w, http.StatusOK, map[string]string{"upload_url": uploadURL})
}

func (h *PublisherHandlers) PublishUploadedWallpaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Validate user ID
	userID, ok := utils.GetStringFromContext(r.Context(), "user_id")
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
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
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
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
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "Wallpaper published successfully",
	})
}

func (h *PublisherHandlers) GetPublishedWallpapers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	publishedWallpapers, err := h.publisherService.GetPublishedWallpapersByUserID(r.Context(), userID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, publishedWallpapers)
}

func (h *PublisherHandlers) GetPublishedWallpapersByDeviceID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	userID, ok := utils.GetStringFromContext(r.Context(), "user_id")
	if !ok {
		utils.WriteJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "missing or invalid user_id",
		})
		return
	}

	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing deviceID",
		})
		return
	}

	publishedWallpapers, err := h.publisherService.
		GetPublishedWallpapersByDeviceID(r.Context(), userID, deviceID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, publishedWallpapers)
}

func (h *PublisherHandlers) ServeWallpaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	deviceID := chi.URLParam(r, "deviceID")
	if deviceID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing deviceID",
		})
		return
	}

	publishedWallpapers, err := h.publisherService.GetPublishedWallpapersByDeviceID(r.Context(), userID, deviceID)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if len(publishedWallpapers) == 0 {
		utils.WriteJSON(w, http.StatusNotFound, map[string]string{
			"error": "no published wallpapers found",
		})
		return
	}

	http.ServeFile(w, r, publishedWallpapers[0].URL)
}

// Delete published wallpaper by hash
func (h *PublisherHandlers) DeletePublishedWallpaperByHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
		return
	}

	// Get the user ID from the request context
	userID := r.Context().Value("user_id").(string)

	hash := chi.URLParam(r, "hash")
	if hash == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "missing hash",
		})
		return
	}

	err := h.publisherService.DeletePublishedWallpaperByHash(r.Context(), userID, hash)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}
}

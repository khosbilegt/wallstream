package handlers

import (
	"net/http"

	"github.io/khosbilegt/wallstream/internal/server/service"
	"github.io/khosbilegt/wallstream/internal/server/utils"
)

type FileHandlers struct {
	fileService *service.FileService
}

func NewFileHandlers(fileService *service.FileService) *FileHandlers {
	return &FileHandlers{fileService: fileService}
}

// Upload wallpaper to the server
func (h *FileHandlers) UploadWallpaper(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	filename, err := h.fileService.UploadFileStream(r.Context(), file, fileHeader.Filename)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"filename": filename})

}

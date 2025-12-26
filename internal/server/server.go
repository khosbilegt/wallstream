package server

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.io/khosbilegt/wallstream/internal/core"
)

// Server handles wallpaper state and file serving
type Server struct {
	states        map[string]*core.PublisherState
	wallpapersDir string
	baseURL       string
	mu            sync.RWMutex
}

// NewServer creates a new server instance
func NewServer(wallpapersDir, baseURL string) (*Server, error) {
	if err := os.MkdirAll(wallpapersDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create wallpapers directory: %w", err)
	}

	return &Server{
		states:        make(map[string]*core.PublisherState),
		wallpapersDir: wallpapersDir,
		baseURL:       baseURL,
	}, nil
}

// GetState returns the current state for a publisher
func (s *Server) GetState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	publisherID := strings.TrimPrefix(r.URL.Path, "/state/")
	if publisherID == "" {
		http.Error(w, "Publisher ID required", http.StatusBadRequest)
		return
	}

	s.mu.RLock()
	state, exists := s.states[publisherID]
	s.mu.RUnlock()

	if !exists {
		http.Error(w, "Publisher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		log.Printf("Failed to encode state: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// UploadWallpaper handles wallpaper uploads from publishers
func (s *Server) UploadWallpaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	publisherID := strings.TrimPrefix(r.URL.Path, "/upload/")
	if publisherID == "" {
		http.Error(w, "Publisher ID required", http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("wallpaper")
	if err != nil {
		http.Error(w, "Failed to get file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Compute hash
	hashBytes := sha256.Sum256(data)
	hash := hex.EncodeToString(hashBytes[:])

	// Determine extension from filename or content type
	ext := "jpg"
	filename := header.Filename
	if filename != "" {
		if dotIdx := strings.LastIndex(filename, "."); dotIdx != -1 {
			ext = strings.ToLower(filename[dotIdx+1:])
		}
	}

	// Save file
	filename = fmt.Sprintf("%s.%s", hash, ext)
	filePath := filepath.Join(s.wallpapersDir, filename)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		log.Printf("Failed to save wallpaper: %v", err)
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Create URL for the wallpaper
	wallpaperURL := fmt.Sprintf("%s/wallpapers/%s", s.baseURL, filename)

	// Update state
	state := &core.PublisherState{
		Hash:      hash,
		URL:       wallpaperURL,
		Timestamp: time.Now().Unix(),
	}

	s.mu.Lock()
	s.states[publisherID] = state
	s.mu.Unlock()

	log.Printf("Uploaded wallpaper for publisher %s: %s", publisherID, hash)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(state); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// ServeWallpaper serves wallpaper files
func (s *Server) ServeWallpaper(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filename := strings.TrimPrefix(r.URL.Path, "/wallpapers/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(s.wallpapersDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Wallpaper not found", http.StatusNotFound)
		return
	}

	// Set content type based on extension
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".bmp":
		w.Header().Set("Content-Type", "image/bmp")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	http.ServeFile(w, r, filePath)
}

// RegisterRoutes registers all HTTP routes
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/state/", s.GetState)
	mux.HandleFunc("/upload/", s.UploadWallpaper)
	mux.HandleFunc("/wallpapers/", s.ServeWallpaper)
}

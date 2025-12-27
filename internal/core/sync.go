package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.io/khosbilegt/wallstream/internal/client/platform"
	"github.io/khosbilegt/wallstream/internal/core/assets"
)

// PublisherState represents the wallpaper state on the server
type PublisherState struct {
	Hash      string `json:"hash"`
	URL       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
}

// Syncer handles fetching the latest publisher state and updating wallpaper
type Syncer struct {
	Config   *Config
	WP       platform.Wallpaper
	LastHash string
	Client   *http.Client
}

// NewSyncer creates a new Syncer instance
func NewSyncer(cfg *Config, wp platform.Wallpaper) *Syncer {
	return &Syncer{
		Config: cfg,
		WP:     wp,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

// PollPublisher polls the server every interval and applies wallpaper if changed
func (s *Syncer) PollPublisher(publisherID string, serverURL string, stopCh <-chan struct{}) {
	ticker := time.NewTicker(s.Config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			err := s.checkAndApply(publisherID, serverURL)
			if err != nil {
				log.Printf("Sync error: %v", err)
			}
		}
	}
}

// checkAndApply fetches the latest state and applies wallpaper if hash changed
func (s *Syncer) checkAndApply(publisherID, serverURL string) error {
	log.Default().Println("Fetching publisher state from", serverURL, "/state/", publisherID)
	url := fmt.Sprintf("%s/state/%s", serverURL, publisherID)
	resp, err := s.Client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch publisher state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var state PublisherState
	if err := json.Unmarshal(body, &state); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	if state.Hash == s.LastHash {
		return nil // no change
	}

	path, err := assets.DownloadImage(state.URL)
	if err != nil {
		return fmt.Errorf("failed to download wallpaper: %w", err)
	}

	if err := s.WP.Set(path); err != nil {
		return fmt.Errorf("failed to set wallpaper: %w", err)
	}

	s.LastHash = state.Hash
	log.Printf("Applied new wallpaper: %s", path)

	return nil
}

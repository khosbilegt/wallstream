package core

import (
	"log"
	"time"

	"github.io/khosbilegt/wallstream/internal/client/platform"
	"github.io/khosbilegt/wallstream/internal/core/assets"
)

// Agent represents the wallpaper agent
type Agent struct {
	Config       *Config
	StateManager *StateManager
	Wallpaper    platform.Wallpaper
	Syncer       *Syncer
	IsPublisher  bool
}

// NewAgent creates a new agent instance
func NewAgent(cfg *Config, wp platform.Wallpaper, isPublisher bool) (*Agent, error) {
	stateMgr, err := NewStateManager(cfg.CacheDir)
	if err != nil {
		return nil, err
	}

	syncer := NewSyncer(cfg, wp)

	return &Agent{
		Config:       cfg,
		StateManager: stateMgr,
		Wallpaper:    wp,
		Syncer:       syncer,
		IsPublisher:  isPublisher,
	}, nil
}

// Run starts the agent loop
func (a *Agent) Run(publisherID, serverURL string, stopCh <-chan struct{}) {
	if a.IsPublisher {
		a.runPublisher(publisherID, serverURL, stopCh)
	} else {
		a.runSubscriber(publisherID, serverURL, stopCh)
	}
}

// runPublisher detects wallpaper changes and pushes to server
func (a *Agent) runPublisher(publisherID, serverURL string, stopCh <-chan struct{}) {
	log.Default().Println("Running publisher")
	state, err := a.StateManager.Load()
	if err != nil {
		log.Printf("Failed to load state: %v", err)
		state = &State{}
	}

	ticker := time.NewTicker(a.Config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			currentPath, err := a.Wallpaper.GetCurrent()
			if err != nil {
				log.Printf("Failed to get current wallpaper: %v", err)
				continue
			}

			hash, err := assets.HashFile(currentPath)
			if err != nil {
				log.Printf("Failed to hash wallpaper: %v", err)
				continue
			}

			if hash == state.LastHash {
				continue // no change
			}

			// Save to cache
			ext := "jpg" // could detect extension dynamically
			cachedPath, err := assets.SaveFileFromPath(currentPath, hash, ext)
			if err != nil {
				log.Printf("Failed to cache wallpaper: %v", err)
				continue
			}

			// TODO: upload to server/CDN and notify server of new hash
			log.Printf("Publisher detected change, wallpaper cached at: %s", cachedPath)

			// Update state
			state.LastHash = hash
			state.Path = cachedPath
			if err := a.StateManager.Save(state); err != nil {
				log.Printf("Failed to save state: %v", err)
			}
		}
	}
}

// runSubscriber polls server and applies wallpaper if changed
func (a *Agent) runSubscriber(publisherID, serverURL string, stopCh <-chan struct{}) {
	log.Default().Println("Running subscriber")
	// Load last state
	state, err := a.StateManager.Load()
	if err != nil {
		log.Printf("Failed to load state: %v", err)
		state = &State{}
	}
	a.Syncer.LastHash = state.LastHash

	a.Syncer.PollPublisher(publisherID, serverURL, stopCh)

	// After poll loop exits, save last state
	state.LastHash = a.Syncer.LastHash
	if err := a.StateManager.Save(state); err != nil {
		log.Printf("Failed to save state: %v", err)
	}
}

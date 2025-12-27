package core

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// State represents the agent’s local wallpaper state
type State struct {
	LastHash  string `json:"last_hash"`
	Path      string `json:"path"`      // local path of applied wallpaper
	Timestamp int64  `json:"timestamp"` // unix timestamp of last applied wallpaper
}

// StateManager handles reading/writing local state to a file
type StateManager struct {
	FilePath string
}

// NewStateManager creates a new StateManager storing state in the given directory
func NewStateManager(dir string) (*StateManager, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create state dir: %w", err)
	}

	return &StateManager{
		FilePath: filepath.Join(dir, "state.json"),
	}, nil
}

// Load reads the state from disk. Returns empty state if file does not exist.
func (m *StateManager) Load() (*State, error) {
	data, err := os.ReadFile(m.FilePath)
	log.Default().Println("Loading state from", m.FilePath)
	if os.IsNotExist(err) {
		log.Default().Println("State file does not exist, creating new state")
		return &State{}, nil
	}
	if err != nil {
		log.Default().Println("Failed to read state file:", err)
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		log.Default().Println("Failed to unmarshal state file, recreating state:", err)
		return &State{}, nil
	}
	return &s, nil
}

// Save writes the current state to disk
func (m *StateManager) Save(s *State) error {
	s.Timestamp = time.Now().Unix()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(m.FilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

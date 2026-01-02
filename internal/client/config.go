package client

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type Config struct {
	CacheDir     string
	PollInterval time.Duration
}

// DefaultConfig returns a cross-platform default configuration
func DefaultConfig() (*Config, error) {
	cacheDir, err := defaultCacheDir()
	if err != nil {
		return nil, err
	}

	return &Config{
		CacheDir:     cacheDir,
		PollInterval: 10 * time.Second,
	}, nil
}

// defaultCacheDir returns OS-specific cache directory
func defaultCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var dir string
	switch runtime.GOOS {
	case "windows":
		local := os.Getenv("LOCALAPPDATA")
		if local == "" {
			return "", fmt.Errorf("LOCALAPPDATA not set")
		}
		dir = filepath.Join(local, "Wallcast", "cache")
	case "darwin":
		dir = filepath.Join(home, "Library", "Caches", "Wallcast")
	case "linux":
		dir = filepath.Join(home, ".cache", "wallcast")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	return dir, nil
}

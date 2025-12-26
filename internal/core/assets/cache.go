package assets

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// DefaultCacheDir returns the OS-specific cache directory for the agent.
func DefaultCacheDir() (string, error) {
	var base string
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "windows":
		base = filepath.Join(os.Getenv("LOCALAPPDATA"), "Wallcast", "cache")
	case "darwin":
		base = filepath.Join(home, "Library", "Caches", "Wallcast")
	case "linux":
		base = filepath.Join(home, ".cache", "wallcast")
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	if err := os.MkdirAll(base, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache dir: %w", err)
	}

	return base, nil
}

// FileExists checks if a file already exists at the given path.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// PathForHash returns the full cache path for a file given its SHA256 hash.
// e.g., abcd1234... -> cache/ab/cd/abcd1234.jpg
func PathForHash(hash string, ext string) (string, error) {
	cacheDir, err := DefaultCacheDir()
	if err != nil {
		return "", err
	}

	// Use first 4 chars of hash as subfolders to avoid too many files in one dir
	if len(hash) < 4 {
		return "", fmt.Errorf("hash too short: %s", hash)
	}

	subDir := filepath.Join(cacheDir, hash[:2], hash[2:4])
	if err := os.MkdirAll(subDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create subdir: %w", err)
	}

	return filepath.Join(subDir, hash+"."+ext), nil
}

// SaveFile writes bytes to the cache path for a given hash.
func SaveFile(data []byte, hash, ext string) (string, error) {
	path, err := PathForHash(hash, ext)
	if err != nil {
		return "", err
	}

	if FileExists(path) {
		return path, nil
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return path, nil
}

func SaveFileFromPath(srcPath, hash, ext string) (string, error) {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("failed to read source file: %w", err)
	}
	return SaveFile(data, hash, ext)
}

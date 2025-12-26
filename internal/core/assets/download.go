package assets

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DownloadImage downloads an image from the given URL and caches it by hash.
// Returns the local cached file path.
func DownloadImage(url string) (string, error) {
	// Step 1: Fetch image
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status code %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return "", fmt.Errorf("URL is not an image: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read image body: %w", err)
	}

	// Step 2: Compute hash
	hashBytes := sha256.Sum256(data)
	hash := fmt.Sprintf("%x", hashBytes[:])

	// Step 3: Determine extension from content-type
	var ext string
	switch contentType {
	case "image/jpeg", "image/jpg":
		ext = "jpg"
	case "image/png":
		ext = "png"
	case "image/bmp":
		ext = "bmp"
	default:
		ext = "bin"
	}

	// Step 4: Save to cache
	path, err := SaveFile(data, hash, ext)
	if err != nil {
		return "", fmt.Errorf("failed to save cached image: %w", err)
	}

	return path, nil
}

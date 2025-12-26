package platform

type Wallpaper interface {
	// GetCurrent returns the absolute path to the current wallpaper.
	// If unsupported, return ErrNotSupported.
	GetCurrent() (string, error)

	// Set sets the wallpaper using an absolute local file path.
	Set(path string) error
}

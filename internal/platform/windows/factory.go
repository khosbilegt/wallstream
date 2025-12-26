//go:build windows

package windows

import "github.io/khosbilegt/wallstream/internal/platform"

func NewWallpaper() platform.Wallpaper {
	return New()
}

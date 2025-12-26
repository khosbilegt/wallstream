//go:build windows

package windows

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"

	"github.io/khosbilegt/wallstream/internal/platform"
)

const (
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02
)

var (
	user32                    = syscall.NewLazyDLL("user32.dll")
	procSystemParametersInfoW = user32.NewProc("SystemParametersInfoW")
)

type Wallpaper struct{}

// New creates a Windows wallpaper controller.
func New() platform.Wallpaper {
	return &Wallpaper{}
}

// Set sets the desktop wallpaper using an absolute local file path.
func (w *Wallpaper) Set(path string) error {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	r1, _, err := procSystemParametersInfoW.Call(
		uintptr(SPI_SETDESKWALLPAPER),
		0,
		uintptr(unsafe.Pointer(ptr)),
		uintptr(SPIF_UPDATEINIFILE|SPIF_SENDCHANGE),
	)

	if r1 == 0 {
		if err != nil && err != syscall.Errno(0) {
			return err
		}
		return fmt.Errorf("SystemParametersInfoW failed")
	}

	return nil
}

// GetCurrent returns the absolute path of the current wallpaper.
func (w *Wallpaper) GetCurrent() (string, error) {
	key, err := registry.OpenKey(
		registry.CURRENT_USER,
		`Control Panel\Desktop`,
		registry.QUERY_VALUE,
	)
	if err != nil {
		return "", err
	}
	defer key.Close()

	path, _, err := key.GetStringValue("WallPaper")
	if err != nil {
		return "", err
	}

	if path == "" {
		return "", platform.ErrNotSupported
	}

	return path, nil
}

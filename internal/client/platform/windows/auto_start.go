package windows

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

// EnableAutoStart adds the program to Windows startup
func EnableAutoStart(appName string) error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot get executable path: %w", err)
	}
	exePath = filepath.Clean(exePath)

	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("cannot open registry key: %w", err)
	}
	defer key.Close()

	if err := key.SetStringValue(appName, exePath); err != nil {
		return fmt.Errorf("cannot set registry value: %w", err)
	}

	return nil
}

// DisableAutoStart removes the program from Windows startup
func DisableAutoStart(appName string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("cannot open registry key: %w", err)
	}
	defer key.Close()

	if err := key.DeleteValue(appName); err != nil {
		return fmt.Errorf("cannot delete registry value: %w", err)
	}

	return nil
}

// IsAutoStartEnabled checks if the program is registered to auto-start
func IsAutoStartEnabled(appName string) (bool, error) {
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
	if err != nil {
		return false, fmt.Errorf("cannot open registry key: %w", err)
	}
	defer key.Close()

	_, _, err = key.GetStringValue(appName)
	if err != nil {
		if err == registry.ErrNotExist {
			return false, nil
		}
		return false, fmt.Errorf("cannot query registry value: %w", err)
	}

	return true, nil
}

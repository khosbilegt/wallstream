package windows

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getlantern/systray"
)

// TrayStatus represents messages from the background sync loop
type TrayStatus struct {
	Message string // e.g., "Idle", "Syncing", "Error"
	Error   error  // optional
}

// RunTray starts the system tray and listens for status updates
func RunTray(statusChan <-chan TrayStatus, quitChan chan struct{}) {
	systray.Run(func() {
		onReady(statusChan, quitChan)
	}, func() {
		// cleanup on exit
	})
}

func onReady(statusChan <-chan TrayStatus, quitChan chan struct{}) {
	// Set initial icon and tooltip
	systray.SetTitle("Wallcast")
	systray.SetTooltip("Wallcast: Idle")

	// Load icon from embedded resource or fallback
	iconData, err := os.ReadFile("icon.ico") // provide your icon here
	if err == nil {
		systray.SetIcon(iconData)
	}

	// Context menu items
	syncItem := systray.AddMenuItem("Sync Now", "Trigger manual sync")
	openConfig := systray.AddMenuItem("Open Config Folder", "Open the configuration folder")
	systray.AddSeparator()
	quitItem := systray.AddMenuItem("Quit", "Exit Wallcast")

	// Goroutine to update tooltip from statusChan
	go func() {
		for {
			select {
			case status := <-statusChan:
				if status.Error != nil {
					systray.SetTooltip(fmt.Sprintf("Wallcast: ERROR - %v", status.Error))
				} else {
					systray.SetTooltip(fmt.Sprintf("Wallcast: %s", status.Message))
				}
			case <-quitChan:
				systray.Quit()
				return
			}
		}
	}()

	// Goroutine for menu actions
	go func() {
		for {
			select {
			case <-syncItem.ClickedCh:
				log.Println("Manual sync triggered")
				// send signal to your sync loop here, e.g., via channel
			case <-openConfig.ClickedCh:
				openConfigFolder()
			case <-quitItem.ClickedCh:
				close(quitChan)
				return
			}
		}
	}()
}

// openConfigFolder opens the config folder for the user
func openConfigFolder() {
	var path string
	if configDir, err := os.UserConfigDir(); err == nil {
		path = filepath.Join(configDir, "Wallcast")
	} else {
		log.Println("Cannot determine user config dir:", err)
		return
	}

	cmd := exec.Command("explorer", path)
	if err := cmd.Start(); err != nil {
		log.Println("Failed to open config folder:", err)
	}
}

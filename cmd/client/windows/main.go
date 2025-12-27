package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.io/khosbilegt/wallstream/internal/client/platform/windows"
)

func main() {
	log.Println("Starting Wallstream Client for Windows...")

	stop := make(chan struct{})
	statusChan := make(chan windows.TrayStatus, 5)
	quitChan := make(chan struct{})

	go func() {
		windows.RunTray(statusChan, quitChan)
	}()

	// Wait for interrupt or quit signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	select {
	case <-c:
		log.Println("Received interrupt signal")
		close(quitChan)
	case <-quitChan:
		log.Println("Received quit signal from tray")
	}

	close(stop)
	log.Println("Wallstream Client stopped")
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.io/khosbilegt/wallstream/internal/client/platform/windows"
	"github.io/khosbilegt/wallstream/internal/core"
)

func main() {
	log.Println("Starting Wallstream Client...")
	cfg, _ := core.DefaultConfig()
	wp := windows.New()

	agent, _ := core.NewAgent(cfg, wp, false)

	stop := make(chan struct{})
	go agent.Run("publisher123", "http://localhost:8080", stop)

	// Wait for Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	close(stop)

	log.Println("Wallstream Client stopped")
}

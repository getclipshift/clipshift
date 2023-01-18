package cmd

import (
	"os"
	"os/signal"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/internal/clip"
	"github.com/jhotmann/clipshift/internal/logger"
	"github.com/jhotmann/clipshift/internal/ui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	log       = logger.Log
	clients   []backends.BackendClient
	interrupt chan os.Signal
)

func init() {
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Keep clipboard in sync",
	Long:  `Subscribes to configured backends and keeps clipboard in sync with other devices`,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(logrus.Fields{
			"Loglevel": config.Logging.Level,
			"Logout":   config.Logging.Destination,
			"Backends": config.Backends,
		}).Info("Sync command")

		// Initialize clients
		for _, b := range config.Backends {
			clients = append(clients, backends.New(b))
		}
		// Handle messages
		for _, client := range clients {
			go client.HandleMessages()
		}

		clip.ClipInit()
		ui.TrayInit()

		interrupt = make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)
		go handleInterrupt()
	},
}

func handleInterrupt() {
	for i := range interrupt {
		log.Info("Interrupt received: ", i.String())
		backends.Close()
		ui.TrayOnExit()
	}
}

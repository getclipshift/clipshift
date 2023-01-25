package cmd

import (
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/internal/clip"
	"github.com/jhotmann/clipshift/internal/ui"
	"github.com/spf13/cobra"
)

var (
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
		errorZeroBackends()
		backendHosts := []string{}
		for _, backend := range config.Backends {
			backendHosts = append(backendHosts, backend.Host)
		}
		log.WithField("Hosts", strings.Join(backendHosts, ", ")).Info("Syncing clipboard")

		// Initialize clients
		for _, b := range config.Backends {
			c := backends.New(b)
			if c != nil {
				clients = append(clients, c)
			}
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
		runtime.Goexit()
	},
}

func handleInterrupt() {
	for i := range interrupt {
		log.Info("Interrupt received: ", i.String())
		backends.Close()
		os.Exit(0)
	}
}

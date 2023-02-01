package cmd

import (
	"fmt"
	"os"

	"github.com/jhotmann/clipshift/backends"
	"github.com/oleiade/reflections"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configAddBackendCmd)
}

var configAddBackendCmd = &cobra.Command{
	Use:   "add-backend",
	Short: "Add a backend server",
	Long:  `Add a backend server to your configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		newBackend := backends.BackendConfig{}

		selectedType, err := pterm.DefaultInteractiveSelect.WithOptions(getBackendHostTypes()).Show()
		if err != nil {
			log.WithError(err).Error("Error getting backend type")
			os.Exit(1)
		}
		newBackend.Type = selectedType

		var hostHint string
		switch selectedType {
		case backends.Hosts.Nostr:
			hostHint = "wss://relay.example.com"
		case backends.Hosts.Ntfy:
			hostHint = "https://ntfy.example.com"
		}

		newBackend.Host = getTextInput("Host (" + hostHint + ")")

		switch selectedType {
		case backends.Hosts.Nostr:
			newBackend.Pass = getTextInput("Nostr private key")
		default:
			newBackend.User = getTextInput("Username")
			newBackend.Pass = getTextInput("Password")
		}

		// platform-specific options
		switch selectedType {
		case backends.Hosts.Ntfy:
			newBackend.Topic = getTextInput("Topic")
			newBackend.EncryptionKey = getTextInput("Encryption key (leave blank to disable encryption)")
		}

		var availableActions []string
		actions, _ := reflections.Items(&backends.SyncActions)
		for _, v := range actions {
			availableActions = append(availableActions, fmt.Sprintf("%v", v))
		}
		selectedAction, err := pterm.DefaultInteractiveSelect.WithOptions(availableActions).Show()
		if err != nil {
			log.WithError(err).Error("Error getting backend action")
			os.Exit(1)
		}
		newBackend.Action = selectedAction

		// Confirm
		printBackendConfig(newBackend)
		confirmation, err := pterm.DefaultInteractiveConfirm.WithDefaultText("Add to configuration").WithDefaultValue(true).Show()
		if err != nil {
			log.WithError(err).Error("Error getting confirmation")
			os.Exit(1)
		}
		if !confirmation {
			os.Exit(0)
		}

		config.Backends = append(config.Backends, newBackend)
		writeConfig()
	},
}

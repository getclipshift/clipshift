package cmd

import (
	"fmt"
	"os"

	"github.com/getclipshift/clipshift/backends"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configAddBackendCmd)
}

var configAddBackendCmd = &cobra.Command{
	Use:     "add-backend",
	Aliases: []string{"add", "a"},
	Short:   "Add a backend server",
	Long:    `Add a backend server to your configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		newBackend := backends.BackendConfig{}

		app = tview.NewApplication()
		box := tview.NewBox().SetBorder(true).SetTitle("Select Backend Type").SetTitleAlign(tview.AlignLeft)
		list := tview.NewList()
		list.Box = box
		list.AddItem("ntfy", "Push notification server", []rune(fmt.Sprintf("%d", 1))[0], nil)
		list.AddItem("nostr", "Encrypted direct messages over a relay", []rune(fmt.Sprintf("%d", 2))[0], nil)
		list.SetSelectedFunc(func(_ int, t string, _ string, _ rune) {
			newBackend.Type = t
			config.Backends = append(config.Backends, newBackend)
			addEditBackendForm(-1)
		})
		if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
			log.WithError(err).Error("Error in TUI")
			os.Exit(1)
		}
	},
}

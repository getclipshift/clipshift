package cmd

import (
	"os"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configEditBackendCmd)
}

var configEditBackendCmd = &cobra.Command{
	Use:     "edit-backend",
	Aliases: []string{"edit", "e"},
	Short:   "Edit a backend",
	Long:    `Edit a backend server from your configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(config.Backends) == 0 {
			pterm.Error.Println("No backends configured, add one with 'clipshift config add-backend'")
			os.Exit(1)
		}
		app = tview.NewApplication()
		if len(args) > 0 {
			if num, err := strconv.ParseInt(args[0], 10, 0); err == nil && int(num) <= len(config.Backends) {
				addEditBackendForm(int(num) - 1)
				os.Exit(0)
			}
		}

		box := tview.NewBox().SetBorder(true).SetTitle("Select Backend").SetTitleAlign(tview.AlignLeft)
		list := tview.NewList()
		list.Box = box
		for i, b := range config.Backends {
			list.AddItem(b.Type, b.Host, indexToRune(i), nil)
		}
		list.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
			addEditBackendForm(i)
		})
		if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
			log.WithError(err).Error("Error in TUI")
			os.Exit(1)
		}
	},
}

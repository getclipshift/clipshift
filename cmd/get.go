package cmd

import (
	"os"
	"strings"

	"github.com/getclipshift/clipshift/backends"
	"github.com/pterm/pterm"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.Flags().Int("backend", 0, "Backend number to get the current clipboard from (starting index = 1, default is all backends)")
	viper.SetDefault("backend", 0)
	viper.BindPFlag("backend", getCmd.Flags().Lookup("backend"))
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get clipboard from one or more backends",
	Long:  `Get clipboard from one or more backends`,
	Run: func(cmd *cobra.Command, args []string) {
		errorZeroBackends()
		specifiedBackend := viper.GetInt("backend")

		if specifiedBackend > len(config.Backends) {
			pterm.Error.Printfln("You specified backend #%d, but only %d are configured", specifiedBackend, len(config.Backends))
			os.Exit(1)
		}

		if specifiedBackend == 0 { // Get from all backends
			for _, b := range config.Backends { // Skip push backends
				if b.Action == backends.SyncActions.Push {
					continue
				}
				b.Action = backends.SyncActions.Pull
				c := backends.New(b)
				if c != nil {
					clients = append(clients, c)
				}
			}
		} else { // Get from specified backend
			b := config.Backends[specifiedBackend-1]
			b.Action = backends.SyncActions.Pull
			c := backends.New(b)
			if c != nil {
				clients = append(clients, c)
			}
		}

		if len(clients) == 0 {
			pterm.Error.Println("No backends to get from")
			os.Exit(1)
		}

		clips := map[string][]string{}
		count := 0
		for _, c := range clients {
			last := c.Get()
			if last != "" {
				hosts, exists := clips[last]
				clips[last] = append(hosts, c.GetConfig().Host)
				if !exists {
					count += 1
				}
			}
		}

		backends.Close()

		switch count {
		case 0:
			pterm.Error.Println("Unable to retrieve any entries")
		case 1:
			for message := range clips {
				backends.ClipReceived(message, "Clipshift TUI")
				if len(message) > 40 {
					message = message[0:37] + "..."
				}
				pterm.Success.Printfln("Clipboard set: %s", message)
			}
		default:
			app = tview.NewApplication()
			box := tview.NewBox().SetBorder(true).SetTitle("Conflict").SetTitleAlign(tview.AlignLeft)
			list := tview.NewList()
			list.Box = box
			index := 0
			for k, v := range clips {
				list.AddItem(k, strings.Join(v, ", "), indexToRune(index), nil)
				index += 1
			}
			list.SetSelectedFunc(func(_ int, message string, _ string, _ rune) {
				app.Stop()
				backends.ClipReceived(message, "Clipshift TUI")
				if len(message) > 40 {
					message = message[0:37] + "..."
				}
				pterm.Success.Printfln("Clipboard set: %s", message)
			})
			if err := app.SetRoot(list, true).EnableMouse(true).Run(); err != nil {
				log.WithError(err).Error("Error in TUI")
				os.Exit(1)
			}
		}
	},
}

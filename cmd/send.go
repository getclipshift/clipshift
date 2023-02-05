package cmd

import (
	"os"

	"github.com/getclipshift/clipshift/backends"
	"github.com/getclipshift/clipshift/internal/clip"
	"github.com/pterm/pterm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().Int("backend", 0, "Backend number to send the current clipboard to (starting index = 1, default is all backends)")
	viper.SetDefault("backend", 0)
	viper.BindPFlag("backend", sendCmd.Flags().Lookup("backend"))

	sendCmd.Flags().Bool("force", false, "Send to all specified backends regardless of configured action")
	viper.SetDefault("force", false)
	viper.BindPFlag("force", sendCmd.Flags().Lookup("force"))
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send clipboard to one or more backends",
	Long:  `Send clipboard to one or more backends`,
	Run: func(cmd *cobra.Command, args []string) {
		errorZeroBackends()
		specifiedBackend := viper.GetInt("backend")
		log.WithFields(logrus.Fields{
			"Loglevel": config.Logging.Level,
			"Logout":   config.Logging.Destination,
			"Backend":  specifiedBackend,
		}).Info("Send command")

		if specifiedBackend > len(config.Backends) {
			pterm.Error.Printfln("You specified backend #%d, but only %d are configured", specifiedBackend, len(config.Backends))
			os.Exit(1)
		}

		if specifiedBackend == 0 { // Send to all backends
			count := len(config.Backends)
			for _, b := range config.Backends {
				if b.Action == backends.SyncActions.Pull && !viper.GetBool("force") { // Skip "action: pull"
					continue
				}
				if count > 1 && b.Action == backends.SyncActions.Manual && !viper.GetBool("force") { // Ask about "action: manual"
					if response, _ := pterm.DefaultInteractiveConfirm.WithDefaultText(b.Host + " is set to manual, are you sure you want to send?").Show(); !response {
						continue
					}
				}
				b.Action = backends.SyncActions.Push
				c := backends.New(b)
				if c != nil {
					clients = append(clients, c)
				}
			}
		} else { // Send to specified backend
			b := config.Backends[specifiedBackend-1]
			if b.Action == backends.SyncActions.Manual { // Assume they meant to send
				b.Action = backends.SyncActions.Push
			} else if b.Action == backends.SyncActions.Pull { // Ask if they meant to send
				if viper.GetBool("force") {
					b.Action = backends.SyncActions.Push
				} else if res, _ := pterm.DefaultInteractiveConfirm.WithDefaultText(b.Host + "is set to pull, are you sure you want to send?").Show(); res {
					b.Action = backends.SyncActions.Push
				}
			}
			c := backends.New(b)
			if c != nil {
				clients = append(clients, c)
			}
		}

		if len(clients) == 0 {
			pterm.Error.Println("No backends to send to")
			os.Exit(1)
		}

		clipContents := clip.Get()

		for _, c := range clients {
			err := c.Post(clipContents)
			if err != nil {
				pterm.Error.Printfln("Error sending clipboard to backend")
				continue
			}
			pterm.Success.Printfln("Clipboard sent")
		}
		defer backends.Close()
	},
}

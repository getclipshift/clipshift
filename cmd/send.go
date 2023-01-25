package cmd

import (
	"os"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/internal/clip"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.Flags().Int("backend", 0, "Backend number to send the current clipboard to (starting index = 1, default is all backends)")
	viper.SetDefault("backend", 0)
	viper.BindPFlag("backend", sendCmd.Flags().Lookup("backend"))
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send clipboard to one or more backends",
	Long:  `Send clipboard to one or more backends`,
	Run: func(cmd *cobra.Command, args []string) {
		errorZeroBackends()
		log.WithFields(logrus.Fields{
			"Loglevel": config.Logging.Level,
			"Logout":   config.Logging.Destination,
			"Backends": config.Backends,
		}).Info("Send command")

		specifiedBackend := viper.GetInt("backend")
		if specifiedBackend > len(config.Backends) {
			log.Errorf("You specified backend #%d, but only %d are configured", specifiedBackend, len(config.Backends))
			os.Exit(1)
		}
		if specifiedBackend == 0 {
			for _, b := range config.Backends {
				c := backends.New(b)
				if c != nil {
					clients = append(clients, c)
				}
			}
		} else {
			c := backends.New(config.Backends[specifiedBackend-1])
			if c != nil {
				clients = append(clients, c)
			}
		}

		for _, c := range clients {
			err := c.Post(clip.Get())
			if err != nil {
				log.WithError(err).Error("Error sending clipboard to backend")
			}
		}
		defer backends.Close()
	},
}

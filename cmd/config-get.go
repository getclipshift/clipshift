package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.AddCommand(configGetCmd)
}

var configGetCmd = &cobra.Command{
	Use:     "get [optional setting name(s)]",
	Aliases: []string{"g"},
	Short:   "Print configuration",
	Long:    `Print all configuration or just the specified value`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 { // print all
			configPrinter("client-name")
			configPrinter("logging.level")
			configPrinter("logging.destination")
			for i := range config.Backends {
				configPrinter(fmt.Sprintf("backends.%d.%s", i, "type"))
				configPrinter(fmt.Sprintf("backends.%d.%s", i, "host"))
				configPrinter(fmt.Sprintf("backends.%d.%s", i, "user"))
				configPrinterSensitive(fmt.Sprintf("backends.%d.%s", i, "pass"))
				configPrinter(fmt.Sprintf("backends.%d.%s", i, "action"))
				configPrinter(fmt.Sprintf("backends.%d.%s", i, "topic"))
				configPrinterSensitive(fmt.Sprintf("backends.%d.%s", i, "encryptionkey"))
			}
		} else {
			for _, arg := range args {
				if strings.HasSuffix(arg, ".pass") {
					configPrinterSensitive(arg)
				} else {
					configPrinter(arg)
				}
			}
		}
	},
}

func configPrinter(key string) {
	val := viper.Get(key)
	if val == nil {
		return
	}
	fmt.Printf("%s: %v\n", key, val)
}

func configPrinterSensitive(key string) {
	val := viper.GetString(key)
	if val == "" {
		return
	}
	half := len(val) / 2
	if half > 10 {
		half = 10
	}
	redacted := val[0:half] + strings.Repeat("*", half)
	fmt.Printf("%s: %s\n", key, redacted)
}

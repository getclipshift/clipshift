package cmd

import (
	"fmt"
	"os"

	"github.com/jhotmann/clipshift/backends"
	"github.com/oleiade/reflections"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `Manage configuration
Use 'get' and 'set' subcommands to get/set properties
Use 'init' subcommand to initialize a new config file
Use 'add-backend' to add a new backend
User 'edit-backend' to edit an existing backend`,
}

func writeConfig() error {
	out, err := yaml.Marshal(config)
	if err != nil {
		log.WithError(err).Println("Error converting config to yaml")
		return err
	}
	err = os.WriteFile(viper.ConfigFileUsed(), out, 0755)
	if err != nil {
		log.WithError(err).Println("Error writing config file")
		return err
	}
	return nil
}

func getTextInput(prompt string) string {
	response, err := pterm.DefaultInteractiveTextInput.WithDefaultText(prompt).Show()
	if err != nil {
		log.WithError(err).Error("Error getting user input")
		os.Exit(1)
	}
	return response
}

func printBackendConfig(b backends.BackendConfig) {
	fieldMap, _ := reflections.Items(&b)
	configText := "Backend Config"
	for k, v := range fieldMap {
		stringVal := fmt.Sprintf("%v", v)
		if stringVal != "" {
			configText = fmt.Sprintf("%s\n  %s: %s", configText, k, stringVal)
		}
	}
	println(configText)
}

func getBackendHostTypes() []string {
	var availableTypes []string
	hosts, _ := reflections.Items(&backends.Hosts)
	for _, v := range hosts {
		availableTypes = append(availableTypes, fmt.Sprintf("%v", v))
	}
	return availableTypes
}

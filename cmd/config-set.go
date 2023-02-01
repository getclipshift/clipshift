package cmd

import (
	"os"
	"regexp"
	"strconv"

	"github.com/oleiade/reflections"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	configCmd.AddCommand(configSetCmd)
}

var configSetCmd = &cobra.Command{
	Use:   "set [setting] [value]",
	Args:  cobra.ExactArgs(2),
	Short: "Set a configuration value",
	Long: `Set a configuration value
Use 'config' with no arguments to print all settings`,
	Run: func(cmd *cobra.Command, args []string) {
		settingName := args[0]
		settingValue := args[1]
		validSettings := regexp.MustCompile(`^(logging\.(level|destination)|client-name|backends\.\d+\.(type|host|user|pass|action|topic|encryptionkey))$`)
		if validSettings.MatchString(args[0]) {
			backendRegex := regexp.MustCompile(`^backends\.(\d+)\.(.+)`)
			if backendRegex.MatchString(settingName) {
				// Backend setting
				backendMatch := backendRegex.FindStringSubmatch(settingName)
				backendNum, _ := strconv.ParseInt(backendMatch[1], 10, 0)
				if int(backendNum) > len(config.Backends) { // Invalid index
					pterm.Error.Printfln(`%d is higher than the number of configured backends (%d)
If you would like to configure a new backend run 'clipshift config add-backend'`, backendNum, len(config.Backends))
					os.Exit(1)
				} else { // Valid index
					log.WithField("Name", backendMatch[2]).Debug("Looking for yaml tag")
					fieldName, _ := reflections.GetFieldNameByTagValue(&config.Backends[backendNum-1], "yaml", backendMatch[2])
					log.WithField("Field Name", fieldName).Debug("Found field")
					if fieldName == "" {
						os.Exit(1)
					}
					log.WithField("Value", settingValue).Debug("Setting value")
					reflections.SetField(&config.Backends[backendNum-1], fieldName, settingValue)
				}
			} else { // Not a backend setting
				switch settingName {
				case "logging.level":
					config.Logging.Level = settingValue
				case "logging.destination":
					config.Logging.Destination = settingValue
				case "client-name":
					config.ClientName = settingValue
				}
			}
			err := writeConfig()
			if err != nil {
				pterm.Error.Printfln("Error writing config: %v", err)
				os.Exit(1)
			}
			pterm.Success.Printfln("%s updated", viper.ConfigFileUsed())
		}
	},
}

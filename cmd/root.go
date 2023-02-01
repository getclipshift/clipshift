package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/internal/logger"
)

var (
	log     = logger.Log
	cfgFile string
	config  Config
)

var rootCmd = &cobra.Command{
	Use:   "clipshift",
	Short: "Clipboard synchronization application",
	Long:  `clipshift - cross-platform clipboard synchronization tool with support for multiple backends`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Error("Error executing command")
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.clipshift/config.yaml)")
	rootCmd.PersistentFlags().String("client-name", "", "Client name: the name of this client to distinguish from other clients")
	rootCmd.PersistentFlags().String("loglevel", "", "Log level: trace, debug, info, warning, or error")
	rootCmd.PersistentFlags().String("logout", "", "Log destination: stdout or a path to a file (doesn't have to exist)")
}

type LogConfig struct {
	Destination string `yaml:"destination"`
	Level       string `yaml:"level"`
}

type Config struct {
	ClientName string                   `yaml:"client-name,omitempty"`
	Logging    LogConfig                `yaml:"logging"`
	Backends   []backends.BackendConfig `yaml:"backends"`
}

func initConfig() {
	if cfgFile != "" { // Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(path.Join(home, ".clipshift"))
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			println("Error reading config file: " + err.Error())
			os.Exit(1)
		}
	}

	hostname, _ := os.Hostname()
	viper.SetDefault("client-name", hostname)
	viper.BindPFlag("client-name", rootCmd.PersistentFlags().Lookup("client-name"))
	viper.BindEnv("client-name", "CLIPSHIFT_CLIENT_NAME")
	viper.SetDefault("logging.level", "error")
	viper.BindPFlag("logging.level", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindEnv("logging.level", "CLIPSHIFT_LOGLEVEL")
	viper.SetDefault("logging.destination", "stdout")
	viper.BindPFlag("logging.destination", rootCmd.PersistentFlags().Lookup("logout"))
	viper.BindEnv("logging.destination", "CLIPSHIFT_LOGOUT")
	viper.SetDefault("backends", []backends.BackendConfig{})

	viper.Unmarshal(&config)

	logger.LoggerInit(viper.GetString("logging.level"), viper.GetString("logging.destination"))
}

func errorZeroBackends() {
	if len(config.Backends) == 0 {
		log.Error("No backends configured")
		os.Exit(1)
	}
}

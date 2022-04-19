package config

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/mitchellh/go-homedir"
)

//go:embed config.example.yml
var exampleConfig []byte

var (
	hostname   string
	configFile string
)

var UserConfig = struct {
	Backend       string `mapstructure:"backend"`
	Host          string `mapstructure:"host"`
	Topic         string `mapstructure:"topic"`
	User          string `mapstructure:"user"`
	Pass          string `mapstructure:"pass"`
	Client        string `mapstructure:"client"`
	LogLevel      string `mapstructure:"loglevel"`
	EncryptionKey string `mapstructure:"encryptionkey"`
}{}

func ConfigInit() {
	hostname, _ = os.Hostname()
	config.AddDriver(yaml.Driver)
	configFile, _ = homedir.Expand("~/.clipshift/config.yml")

	err := config.LoadFiles(configFile)
	if err != nil {
		fmt.Println("ERROR: config file not found, default config created")
		os.MkdirAll(strings.Replace(configFile, "config.yml", "", 1), 0600)
		os.WriteFile(configFile, exampleConfig, 0600)
		os.Exit(1)
	}
	err = config.BindStruct("", &UserConfig)
	if err != nil {
		panic(err)
	}
	if UserConfig.Client == "" {
		UserConfig.Client = hostname
	}
}

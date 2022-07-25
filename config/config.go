package config

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/jhotmann/clipshift/logger"
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
		logger.Log.Error("Config file not found, default config created")
		os.MkdirAll(strings.Replace(configFile, "config.yml", "", 1), 0600)
		os.WriteFile(configFile, exampleConfig, 0600)
		openConfig(configFile)
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

func openConfig(filePath string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		cmd = exec.Command("cmd", "/C", "start", filepath.FromSlash(filePath))
	case "darwin":
		cmd = exec.Command("open", filePath)
	}
	err := cmd.Start()
	if err == nil {
		cmd.Wait()
	} else {
		logger.Log.WithField("Error", err.Error()).Error("Couldn't open config in editor")
	}

}

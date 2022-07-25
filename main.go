package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/config"
	"github.com/jhotmann/clipshift/logger"
	"github.com/jhotmann/clipshift/ui"
)

var (
	hostname  string
	interrupt chan os.Signal
)

func main() {
	config.ConfigInit()
	logger.LoggerInit(config.UserConfig.LogLevel)

	interrupt = make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go handleInterrupt()

	defer func() {
		if err := recover(); err != nil {
			logger.Log.Error("Exception:", err)
			backends.Close()
			main()
			return
		}
	}()

	backends.BackendInit()
	ClipInit()
	ui.TrayInit()
}

func handleInterrupt() {
	for i := range interrupt {
		fmt.Println("Interrupt received: ", i.String())
		backends.Close()
		ui.TrayOnExit()
	}
}

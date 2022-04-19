package ui

import (
	_ "embed"
	"os"
	"runtime"

	"github.com/getlantern/systray"
)

//go:embed clipboard.png
var iconPng []byte

//go:embed clipboard.ico
var iconIco []byte

func TrayInit() {
	systray.Run(trayOnReady, TrayOnExit)
}

func trayOnReady() {
	systray.SetTitle("")

	if runtime.GOOS == "windows" {
		systray.SetIcon(iconIco)
	} else {
		systray.SetIcon(iconPng)
	}

	mQuit := systray.AddMenuItem("Exit", "Exit clipshift")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func TrayOnExit() {
	os.Exit(0)
}

func TraySetTooltip(msg string) {
	systray.SetTooltip(msg)
}

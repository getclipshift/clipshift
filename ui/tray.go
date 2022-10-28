package ui

import (
	_ "embed"
	"os"
	"runtime"

	"github.com/getlantern/systray"
	"github.com/mitchellh/go-homedir"
)

var (
	macPlistPath string
	startup      bool
)

//go:embed clipboard.png
var iconPng []byte

//go:embed clipboard.ico
var iconIco []byte

func TrayInit() {
	macPlistPath, _ = homedir.Expand("~/Library/LaunchAgents/io.github.clipshift.plist")
	startup = getLaunchAtStartup()
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
	mStartup := systray.AddMenuItemCheckbox("Run at startup", "Run at startup", startup)

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go func() {
		<-mStartup.ClickedCh
		startup = !startup
		if startup {
			mStartup.Check()
		} else {
			mStartup.Uncheck()
		}
		setLaunchAtStartup(startup)
	}()
}

func TrayOnExit() {
	os.Exit(0)
}

func TraySetTooltip(msg string) {
	systray.SetTooltip(msg)
}

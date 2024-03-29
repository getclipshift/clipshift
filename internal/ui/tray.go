//go:build tray

package ui

import (
	_ "embed"
	"os"
	"runtime"

	"github.com/getclipshift/clipshift/backends"
	"github.com/getlantern/systray"
)

//go:embed clipboard.png
var iconPng []byte

//go:embed clipboard.ico
var iconIco []byte

func init() {
	trayEnabled = true
}

func trayRun() {
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
		for {
			<-mStartup.ClickedCh
			startup = !startup
			if startup {
				mStartup.Check()
			} else {
				mStartup.Uncheck()
			}
			setLaunchAtStartup(startup)
		}
	}()
}

func TrayOnExit() {
	backends.Close()
	os.Exit(0)
}

func TraySetTooltip(msg string) {
	systray.SetTooltip(msg)
}

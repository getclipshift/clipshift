package ui

import (
	"runtime"

	"github.com/mitchellh/go-homedir"
)

var (
	macPlistPath string
	winLnkPath   string
	startup      bool
	trayEnabled  bool
)

func TrayInit() {
	if trayEnabled {
		switch runtime.GOOS {
		case "darwin":
			macPlistPath, _ = homedir.Expand("~/Library/LaunchAgents/io.github.clipshift.plist")
		case "windows":
			winLnkPath, _ = homedir.Expand("~/AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup/clipshift.lnk")
		}
		startup = getLaunchAtStartup()
		trayRun()
	}
}

package ui

import (
	"runtime"

	"github.com/mitchellh/go-homedir"
)

var (
	macPlistPath string
	startup      bool
	trayEnabled  bool
)

func TrayInit() {
	if trayEnabled {
		if runtime.GOOS == "darwin" {
			macPlistPath, _ = homedir.Expand("~/Library/LaunchAgents/io.github.clipshift.plist")
		}
		startup = getLaunchAtStartup()
		trayRun()
	}
}

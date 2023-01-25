//go:build !tray

package ui

func trayRun() {
	// This is just to get rid of some warnings of unused items due to build tags
	if trayEnabled {
		setLaunchAtStartup(startup)
	}
}

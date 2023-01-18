package ui

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"text/template"

	"github.com/jhotmann/clipshift/internal/logger"
	"github.com/mitchellh/go-homedir"
)

func getLaunchAtStartup() bool {
	switch runtime.GOOS {
	case "darwin":
		return macGetLaunchAtStartup()
	default:
		return false
	}
}

func macGetLaunchAtStartup() bool {
	if fileExists(macPlistPath) {
		data, err := os.ReadFile(macPlistPath)
		if err != nil {
			logger.Log.WithError(err).Error("Error reading launch agent file")
			return false
		}
		return strings.Contains(string(data), "<key>RunAtLoad</key>\n\t\t<true/>")
	} else {
		macSetLaunchAtStartup(false)
		return false
	}
}

func setLaunchAtStartup(enabled bool) {
	switch runtime.GOOS {
	case "darwin":
		macSetLaunchAtStartup(enabled)
	}
}

type plistData struct {
	Path    string
	Enabled bool
}

func macSetLaunchAtStartup(enabled bool) {
	plistTemplate := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
		<key>Label</key>
		<string>io.github.clipshift</string>
		<key>ProgramArguments</key>
		<array>
				<string>{{ .Path }}</string>
		</array>
		<key>KeepAlive</key>
		<false/>
		<key>RunAtLoad</key>
		<{{ .Enabled }}/>
</dict>
</plist>`
	path, _ := os.Executable()
	data := plistData{
		Path:    path,
		Enabled: enabled,
	}
	tpl := template.Must(template.New("plist").Parse(plistTemplate))
	var out bytes.Buffer
	tpl.Execute(&out, data)
	pListFile, _ := homedir.Expand(macPlistPath)
	os.WriteFile(pListFile, out.Bytes(), 0655)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

package ui

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/getclipshift/clipshift/internal/logger"
	"github.com/mitchellh/go-homedir"
)

func getLaunchAtStartup() bool {
	switch runtime.GOOS {
	case "darwin":
		return macGetLaunchAtStartup()
	case "windows":
		return fileExists(winLnkPath)
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
	case "windows":
		winSetLaunchAtStartup(enabled)
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
	err := os.WriteFile(pListFile, out.Bytes(), 0655)
	if err != nil {
		logger.Log.WithError(err).Error("Error setting startup status")
		return
	}
	logger.Log.WithField("enabled", enabled).Info("Startup status set")
}

func winSetLaunchAtStartup(enabled bool) {
	if enabled {
		exePath, _ := os.Executable()
		exeDir := filepath.Dir(exePath)
		lnkAbs, _ := filepath.Abs(winLnkPath)
		_, err := exec.Command("powershell", "-nologo", "-noprofile", fmt.Sprintf("$s=(New-Object -COM WScript.Shell).CreateShortcut('%s');$s.TargetPath='cmd.exe';$s.Arguments='/c \"start clipshift-tray.exe\"';$s.WorkingDirectory='%s';$s.Save()", lnkAbs, exeDir)).CombinedOutput()
		if err != nil {
			logger.Log.WithError(err).Error("Error creating startup link")
		}
	} else {
		os.Remove(winLnkPath)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

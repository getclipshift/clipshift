package clip

import (
	"context"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/internal/logger"
	"golang.design/x/clipboard"
)

var (
	clipTextChannel <-chan []byte
)

func Get() string {
	contents := clipboard.Read(clipboard.FmtText)
	return string(contents)
}

func ClipInit() {
	clipTextChannel = clipboard.Watch(context.TODO(), clipboard.FmtText)
	go clipTextMonitor()
	logger.Log.Debug("Clipboard is being monitored")
}

func clipTextMonitor() {
	for data := range clipTextChannel {
		clip := string(data)
		if clip != backends.LastReceived {
			logger.Log.WithField("Clipboard", clip).Debug("Clipboard set")
			backends.PostClip(clip)
		}
	}
}

package main

import (
	"context"

	"github.com/jhotmann/clipshift/backends"
	"github.com/jhotmann/clipshift/logger"
	"golang.design/x/clipboard"
)

var (
	clipTextChannel <-chan []byte
)

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

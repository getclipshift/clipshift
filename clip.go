package main

import (
	"context"
	"fmt"

	"github.com/jhotmann/clipshift/backends"
	"golang.design/x/clipboard"
)

var (
	clipTextChannel <-chan []byte
)

func ClipInit() {
	clipTextChannel = clipboard.Watch(context.TODO(), clipboard.FmtText)
	go clipTextMonitor()
	println("Clipboard is being monitored")
}

func clipTextMonitor() {
	for data := range clipTextChannel {
		clip := string(data)
		if clip != LastReceived {
			fmt.Println("Clipboard set: ", clip)
			backends.PostClip(clip)
		}
	}
}

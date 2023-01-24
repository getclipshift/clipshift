package backends

import (
	"github.com/jhotmann/clipshift/internal/logger"
	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

var (
	log          = logger.Log
	LastReceived string

	clients     []BackendClient
	SyncActions = struct {
		Push string
		Pull string
		Sync string
	}{
		Push: "push",
		Pull: "pull",
		Sync: "sync",
	}
)

type BackendConfig struct {
	Type          string `yaml:"type"`
	Host          string `yaml:"host"`
	User          string `yaml:"user"`
	Pass          string `yaml:"pass"`
	Topic         string `yaml:"topic"`
	EncryptionKey string `yaml:"encryptionkey"`
	Action        string `yaml:"action"`
}

type BackendClient interface {
	HandleMessages()
	Post(string) error
	Close()
}

func New(config BackendConfig) BackendClient {
	var client BackendClient
	switch config.Type {
	case "nostr":
		client = nostrInitialize(config)
	case "ntfy":
		client = ntfyInitialize(config)
	}
	clients = append(clients, client)
	return client
}

func Close() {
	for _, c := range clients {
		c.Close()
	}
}

func PostClip(clip string) {
	for _, c := range clients {
		c.Post(clip)
	}
}

func ClipReceived(clip string, client string) {
	if clip == LastReceived {
		return
	}
	LastReceived = clip

	clipboard.Write(clipboard.FmtText, []byte(LastReceived))
	log.WithFields(logrus.Fields{
		"Client":  client,
		"Content": LastReceived,
	}).Debug("Clipboard received")
}

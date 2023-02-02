package backends

import (
	"strings"

	"github.com/jhotmann/clipshift/internal/logger"
	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

var (
	log          = logger.Log
	LastReceived string

	clients     []BackendClient
	SyncActions = struct {
		Push   string
		Pull   string
		Sync   string
		Manual string
	}{
		Push:   "push",
		Pull:   "pull",
		Sync:   "sync",
		Manual: "manual",
	}
	Hosts = struct {
		Ntfy  string
		Nostr string
	}{
		Ntfy:  "ntfy",
		Nostr: "nostr",
	}
)

type BackendConfig struct {
	Type          string `yaml:"type"`
	Host          string `yaml:"host"`
	User          string `yaml:"user,omitempty"`
	Pass          string `yaml:"pass,omitempty"`
	Topic         string `yaml:"topic,omitempty"`
	EncryptionKey string `yaml:"encryptionkey,omitempty"`
	Action        string `yaml:"action"`
	Compression   bool   `yaml:"compression,omitempty"`
}

type BackendClient interface {
	HandleMessages()
	Post(string) error
	Close()
}

func New(config BackendConfig) BackendClient {
	var client BackendClient
	switch strings.ToLower(config.Type) {
	case Hosts.Nostr:
		client = nostrInitialize(config)
	case Hosts.Ntfy:
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

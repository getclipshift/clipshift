package backends

import (
	"github.com/jhotmann/clipshift/internal/logger"
	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

var (
	log          = logger.Log
	LastReceived string
	// nonce        []byte
	// aead         cipher.AEAD

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

// func BackendInit() {
// 	configuredBackend = "" //config.UserConfig.Backend

// 	switch configuredBackend {
// 	case "ntfy":
// 		ntfyInit()
// 	case "nostr":
// 		// nostrInit()
// 	default:
// 		logger.Log.Fatalf("Invalid backend '%s'", configuredBackend)
// 	}

// 	encryptionEnabled = false //config.UserConfig.EncryptionKey != ""
// 	if encryptionEnabled {
// 		encryptionkey = sha256.Sum256([]byte("")) //config.UserConfig.EncryptionKey))
// 		nonce = make([]byte, chacha20poly1305.NonceSizeX)
// 		aead, _ = chacha20poly1305.NewX(encryptionkey[:])
// 	}
// }

func Close() {
	for _, c := range clients {
		c.Close()
	}
}

func PostClip(clip string) {
	for _, c := range clients {
		c.Post(clip)
	}
	// if encryptionEnabled {
	// 	old := clip
	// 	clip = base64.StdEncoding.EncodeToString(encryptString(clip))
	// 	logger.Log.WithFields(logrus.Fields{
	// 		"Old": old,
	// 		"New": clip,
	// 	}).Info("Encrypted clip")
	// }

	// switch configuredBackend {
	// case "ntfy":
	// 	ntfyPostClip(clip)
	// case "nostr":
	// 	// nostrPostClip(clip)
	// }
}

func ClipReceived(clip string, client string) {
	// if encryptionEnabled {
	// 	lastBytes, _ := base64.StdEncoding.DecodeString(clip)
	// 	clip = decryptBytes(lastBytes)
	// }

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

// func encryptString(msg string) []byte {
// 	return aead.Seal(nil, nonce, []byte(msg), nil)
// }

// func decryptBytes(cipher []byte) string {
// 	decrypted, _ := aead.Open(nil, nonce, []byte(cipher), nil)
// 	return string(decrypted)
// }

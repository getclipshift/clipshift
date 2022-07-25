package backends

import (
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/jhotmann/clipshift/config"
	"github.com/jhotmann/clipshift/logger"
	"github.com/jhotmann/clipshift/ui"
	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	LastReceived      string
	configuredBackend string
	encryptionEnabled bool
	encryptionkey     [32]byte
	nonce             []byte
	aead              cipher.AEAD
)

func BackendInit() {
	configuredBackend = config.UserConfig.Backend

	switch configuredBackend {
	case "ntfy":
		ntfyInit()
		break
	default:
		// TODO invalid backend
	}

	encryptionEnabled = config.UserConfig.EncryptionKey != ""
	if encryptionEnabled {
		encryptionkey = sha256.Sum256([]byte(config.UserConfig.EncryptionKey))
		nonce = make([]byte, chacha20poly1305.NonceSizeX)
		aead, _ = chacha20poly1305.NewX(encryptionkey[:])
	}
}

func Close() {
	switch configuredBackend {
	case "ntfy":
		ntfyStreamClose()
		break
	}
}

func PostClip(clip string) {
	if encryptionEnabled {
		old := clip
		clip = base64.StdEncoding.EncodeToString(encryptString(clip))
		logger.Log.WithFields(logrus.Fields{
			"Old": old,
			"New": clip,
		}).Info("Encrypted clip")
	}

	switch configuredBackend {
	case "ntfy":
		ntfyPostClip(clip)
		break
	}
}

func ClipReceived(clip string, client string) {
	if encryptionEnabled {
		lastBytes, _ := base64.StdEncoding.DecodeString(clip)
		clip = decryptBytes(lastBytes)
	}

	if client == config.UserConfig.Client || clip == LastReceived {
		return
	}
	LastReceived = clip

	clipboard.Write(clipboard.FmtText, []byte(LastReceived))
	ui.TraySetTooltip(fmt.Sprintf("%s - %s", time.Now().Format("20060102 15:04:05"), client))
	logger.Log.WithFields(logrus.Fields{
		"Client":  client,
		"Content": LastReceived,
	}).Debug("Clipboard received")
}

func encryptString(msg string) []byte {
	return aead.Seal(nil, nonce, []byte(msg), nil)
}

func decryptBytes(cipher []byte) string {
	decrypted, _ := aead.Open(nil, nonce, []byte(cipher), nil)
	return string(decrypted)
}

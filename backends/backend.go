package backends

import (
	"crypto/cipher"
	"crypto/sha256"

	"github.com/jhotmann/clipshift/config"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	lastReceived      string
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
		clip = encryptString(clip)
	}

	switch configuredBackend {
	case "ntfy":
		ntfyPostClip(clip)
		break
	}
}

func encryptString(msg string) string {
	return string(aead.Seal(nil, nonce, []byte(msg), nil))
}

func decryptString(cipher string) string {
	decrypted, _ := aead.Open(nil, nonce, []byte(cipher), nil)
	return string(decrypted)
}

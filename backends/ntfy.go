package backends

import (
	"encoding/base64"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ntfyClient "heckel.io/ntfy/client"
)

type NtfyClient struct {
	Config        BackendConfig
	ClientName    string
	Client        *ntfyClient.Client
	EncryptionKey string
}

func ntfyInitialize(config BackendConfig) *NtfyClient {
	c := NtfyClient{
		Config:        config,
		ClientName:    viper.GetString("client-name"),
		EncryptionKey: config.EncryptionKey,
	}
	c.Client = ntfyClient.New(&ntfyClient.Config{
		DefaultHost: config.Host,
	})
	log.WithFields(logrus.Fields{
		"Host":  config.Host,
		"User":  config.User,
		"Topic": config.Topic,
	}).Info("Connecting to ntfy relay")
	c.Client.Subscribe(config.Topic, ntfyClient.WithBasicAuth(config.User, config.Pass))
	return &c
}

func (c *NtfyClient) HandleMessages() {
	if c.Config.Action == SyncActions.Push {
		return
	}
	for m := range c.Client.Messages {
		log.WithFields(logrus.Fields{
			"Title":   m.Title,
			"Message": m.Message,
		}).Debug("Ntfy message received")
		if m.Title == c.ClientName {
			continue
		}
		if c.EncryptionKey != "" {
			lastBytes, _ := base64.StdEncoding.DecodeString(m.Message)
			m.Message = decryptBytes(lastBytes)
		}
		ClipReceived(m.Message, m.Title)
	}
}

func (c *NtfyClient) Post(clip string) error {
	if c.Config.Action == SyncActions.Pull {
		return nil
	}
	if c.EncryptionKey != "" {
		clip = base64.StdEncoding.EncodeToString(encryptString(clip))
	}
	_, err := c.Client.Publish(c.Config.Topic, clip, ntfyClient.WithTitle(c.ClientName), ntfyClient.WithBasicAuth(c.Config.User, c.Config.Pass), ntfyClient.WithPriority("1"))
	if err != nil {
		log.WithError(err).Errorf("Error sending clipboard to ntfy host %s", c.Config.Host)
	} else {
		log.Debugf("Clipboard sent to ntfy host %s", c.Config.Host)
	}
	return err
}

func (c *NtfyClient) Close() {
	log.Debug("Closing ntfy subscription")
	c.Client.UnsubscribeAll(c.Config.Topic)
}

package backends

import (
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
		ClipReceived(m.Message, m.Title)
	}
}

func (c *NtfyClient) Post(clip string) error {
	c.Client.Publish(c.Config.Topic, clip, ntfyClient.WithTitle(c.ClientName), ntfyClient.WithBasicAuth(c.Config.User, c.Config.Pass), ntfyClient.WithPriority("1"))
	return nil
}

func (c *NtfyClient) Close() {
	log.Debug("Closing ntfy subscription")
	c.Client.UnsubscribeAll(c.Config.Topic)
}

package backends

import (
	"github.com/jhotmann/clipshift/internal/aes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ntfyClient "heckel.io/ntfy/client"
)

type NtfyClient struct {
	Config        BackendConfig
	ClientName    string
	Client        *ntfyClient.Client
	EncryptionKey []byte
	BaseOptions   []ntfyClient.PublishOption
}

func ntfyInitialize(config BackendConfig) *NtfyClient {
	c := NtfyClient{
		Config:     config,
		ClientName: viper.GetString("client-name"),
	}
	if config.EncryptionKey != "" {
		c.EncryptionKey = aes.GetKey(config.EncryptionKey)
	}
	c.Client = ntfyClient.New(&ntfyClient.Config{
		DefaultHost: config.Host,
	})
	log.WithFields(logrus.Fields{
		"Host":  config.Host,
		"User":  config.User,
		"Topic": config.Topic,
	}).Info("Connecting to ntfy relay")
	if config.User != "" && config.Pass != "" {
		c.BaseOptions = []ntfyClient.PublishOption{ntfyClient.WithBasicAuth(config.User, config.Pass)}
	}
	c.Client.Subscribe(config.Topic, c.BaseOptions...)
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
		if c.EncryptionKey != nil {
			var err error
			m.Message, err = aes.Decrypt(c.EncryptionKey, m.Message)
			if err != nil {
				log.WithError(err).Error("Error decrypting clipboard")
				return
			}
		}
		ClipReceived(m.Message, m.Title)
	}
}

func (c *NtfyClient) Post(clip string) error {
	if c.Config.Action == SyncActions.Pull {
		return nil
	}
	if c.EncryptionKey != nil {
		var err error
		clip, err = aes.Encrypt(c.EncryptionKey, clip)
		if err != nil {
			log.WithError(err).Error("Error encrypting clipboard")
			return err
		}
	}
	opts := append(c.BaseOptions, ntfyClient.WithTitle(c.ClientName), ntfyClient.WithPriority("1"))
	_, err := c.Client.Publish(c.Config.Topic, clip, opts...)
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

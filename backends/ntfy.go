package backends

import (
	"time"

	"github.com/getclipshift/clipshift/internal/aes"
	"github.com/golang-module/dongle"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ntfyClient "heckel.io/ntfy/client"
	ntfyLog "heckel.io/ntfy/log"
)

type NtfyClient struct {
	Config      BackendConfig
	ClientName  string
	Client      *ntfyClient.Client
	Cipher      *dongle.Cipher
	BaseOptions []ntfyClient.PublishOption
}

func ntfyInitialize(config BackendConfig) *NtfyClient {
	c := NtfyClient{
		Config:     config,
		ClientName: viper.GetString("client-name"),
	}
	if config.EncryptionKey != "" {
		c.Cipher = aes.GetCypher(config.EncryptionKey)
	}
	ntfyLog.SetLevel(ntfyLog.ErrorLevel)
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
	if c.Config.Action == SyncActions.Push || c.Config.Action == SyncActions.Manual {
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
		if c.Cipher != nil {
			m.Message = aes.Decrypt(c.Cipher, m.Message)
		}
		ClipReceived(m.Message, m.Title)
	}
}

func (c *NtfyClient) Post(clip string) error {
	if c.Config.Action == SyncActions.Pull || c.Config.Action == SyncActions.Manual {
		return nil
	}
	if c.Cipher != nil {
		clip = aes.Encrypt(c.Cipher, clip)
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

func (c *NtfyClient) Get() string {
	opts := append(c.BaseOptions, ntfyClient.WithSinceDuration(24*time.Hour))
	messages, err := c.Client.Poll(c.Config.Topic, opts...)
	if err != nil {
		log.WithError(err).Error("Error getting ntfy messages")
		return ""
	}
	log.WithFields(logrus.Fields{
		"Count": len(messages),
	}).Info("Queried Ntfy messages")
	if len(messages) == 0 {
		return ""
	}
	message := messages[len(messages)-1].Message
	if c.Config.EncryptionKey != "" {
		message = aes.Decrypt(c.Cipher, message)
	}
	return message
}

func (c *NtfyClient) GetConfig() *BackendConfig {
	return &c.Config
}

func (c *NtfyClient) Close() {
	log.Debug("Closing ntfy subscription")
	c.Client.UnsubscribeAll(c.Config.Topic)
}

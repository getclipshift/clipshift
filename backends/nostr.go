package backends

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/getclipshift/clipshift/internal/logger"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type NostrClient struct {
	Config       BackendConfig
	Client       string
	Ctx          context.Context
	Cancel       context.CancelFunc
	Relay        *nostr.Relay
	Subscription *nostr.Subscription
	SharedSecret []byte
}

func nostrInitialize(config BackendConfig) *NostrClient {
	var err error
	if config.Pass != "" && config.User == "" {
		config.User, err = nostr.GetPublicKey(config.Pass)
		if err != nil {
			log.WithError(err).Error("Invalid password configured")
			return nil
		}
	}
	c := NostrClient{
		Config: config,
		Client: viper.GetString("client-name"),
	}
	secret, err := nip04.ComputeSharedSecret(config.User, config.Pass)
	if err != nil {
		log.WithError(err).Fatal("Unable to compute shared secret")
		return &c
	}
	c.SharedSecret = secret
	c.Ctx, c.Cancel = context.WithCancel(context.Background())
	log.WithFields(logrus.Fields{
		"Host": config.Host,
		"User": config.User,
	}).Info("Connecting to nostr relay")
	c.Relay, err = nostr.RelayConnect(c.Ctx, config.Host)
	if err != nil {
		log.WithError(err).Fatal("Unable to connect to configured relay")
	}
	filter := nostr.Filter{
		Kinds:   []int{4},
		Authors: []string{config.User},
		Limit:   1,
	}
	c.Subscription = c.Relay.Subscribe(c.Ctx, nostr.Filters{filter})

	go func() {
		<-c.Subscription.EndOfStoredEvents
		// TODO - should I do anything with this?
	}()

	return &c
}

func (c *NostrClient) HandleMessages() {
	if c.Config.Action == SyncActions.Push || c.Config.Action == SyncActions.Manual {
		return
	}
	for ev := range c.Subscription.Events {
		log.WithField("Event", ev).Debug("Nostr message received")
		clientName, message := decryptNostrMessage(ev.Content, c.SharedSecret)
		if clientName == "" && message == "" {
			continue
		} else if clientName == c.Client {
			continue
		}
		ClipReceived(message, clientName)
	}
}

func (c *NostrClient) Post(clip string) error {
	if c.Config.Action == SyncActions.Pull || c.Config.Action == SyncActions.Manual {
		log.WithField("Action", c.Config.Action).Debug("Not posting clipboard due to configured Action")
		return nil
	}
	encrypted, err := nip04.Encrypt(fmt.Sprintf("%s---%s", c.Client, clip), c.SharedSecret)
	if err != nil {
		logger.Log.WithError(err).Error("Unable to encrypt clipboard")
		return err
	}
	event := nostr.Event{
		PubKey:    c.Config.User,
		CreatedAt: time.Now(),
		Kind:      4,
		Content:   encrypted,
		Tags: nostr.Tags{
			nostr.Tag{"p", c.Config.User},
		},
	}
	event.Sign(c.Config.Pass)
	status := c.Relay.Publish(c.Ctx, event)
	if status.String() == "failed" {
		return fmt.Errorf("%s received from relay", status)
	}
	log.Debugf("Clipboard sent to nostr relay %s", c.Config.Host)
	return nil
}

func (c *NostrClient) Get() string {
	if c == nil {
		return ""
	}
	pastDay := time.Now().Add(-24 * time.Hour)
	filter := nostr.Filter{
		Kinds:   []int{4},
		Authors: []string{c.Config.User},
		Limit:   1,
		Since:   &pastDay,
	}
	c.Relay.Connect(c.Ctx)
	events := c.Relay.QuerySync(c.Ctx, filter)
	log.WithFields(logrus.Fields{
		"Count": len(events),
	}).Debug("Queried Nostr events")
	if len(events) == 0 {
		return ""
	}
	clientName, message := decryptNostrMessage(events[0].Content, c.SharedSecret)
	if clientName == "" && message == "" {
		return ""
	}
	return message
}

func (c *NostrClient) GetConfig() *BackendConfig {
	return &c.Config
}

func (c *NostrClient) Close() {
	log.Debug("Closing nostr stream")
	c.Cancel()
	c.Relay.Connection.Close()
	c.Relay.Close()
}

func decryptNostrMessage(encrypted string, secret []byte) (clientName string, message string) {
	content, err := nip04.Decrypt(encrypted, secret)
	if err != nil {
		logger.Log.WithError(err).Error("Error decrypting Nostr message contnet")
		return "", ""
	}
	parts := strings.SplitN(content, "---", 2)
	if len(parts) != 2 {
		log.WithField("Content", content).Error("Ignoring Nostr message because it is an incorrect format")
		return "", ""
	}
	return parts[0], parts[1]
}

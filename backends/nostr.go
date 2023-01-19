package backends

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jhotmann/clipshift/internal/logger"
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
	c := NostrClient{
		Config: config,
		Client: viper.GetString("client-name"),
	}
	secret, err := nip04.ComputeSharedSecret(config.Pass, config.User)
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
	if c.Config.Action == SyncActions.Push {
		return
	}
	for ev := range c.Subscription.Events {
		log.WithField("Event", ev).Debug("Nostr message received")
		content, err := nip04.Decrypt(ev.Content, c.SharedSecret)
		if err != nil {
			logger.Log.WithError(err).Error("Error decrypting Nostr message contnet")
			continue
		}
		parts := strings.SplitN(content, "---", 2)
		if len(parts) != 2 {
			log.WithField("Content", content).Error("Ignoring Nostr message because it is an incorrect format")
			continue
		}
		if parts[0] == c.Client {
			continue
		}
		ClipReceived(parts[1], parts[0])
	}
}

func (c *NostrClient) Post(clip string) error {
	if c.Config.Action == SyncActions.Pull {
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
	}
	event.Sign(c.Config.Pass)
	status := c.Relay.Publish(c.Ctx, event)
	if status.String() == "failed" {
		return fmt.Errorf("%s received from relay", status)
	}
	log.Debugf("Clipboard sent to nostr relay %s", c.Config.Host)
	return nil
}

func (c *NostrClient) Close() {
	log.Debug("Closing nostr stream")
	c.Cancel()
	c.Relay.Connection.Close()
	c.Relay.Close()
}

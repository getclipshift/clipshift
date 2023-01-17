package backends

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jhotmann/clipshift/config"
	"github.com/jhotmann/clipshift/logger"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/sirupsen/logrus"
)

var (
	nostrContext       context.Context
	nostrContextCancel context.CancelFunc
	nostrRelay         *nostr.Relay
	nostrSub           *nostr.Subscription
	nostrSharedSecret  []byte
)

func nostrInit() {
	var err error
	nostrSharedSecret, err = nip04.ComputeSharedSecret(config.UserConfig.Pass, config.UserConfig.User)
	if err != nil {
		logger.Log.WithError(err).Fatal("Unable to compute shared secret")
	}
	logger.Log.WithFields(logrus.Fields{
		"Host": config.UserConfig.Host,
		"User": config.UserConfig.User,
	}).Info("Connecting to nostr relay")
	nostrContext, nostrContextCancel = context.WithCancel(context.Background())
	nostrRelay, err = nostr.RelayConnect(nostrContext, config.UserConfig.Host)
	if err != nil {
		logger.Log.WithError(err).Fatal("Unable to connect to configured relay")
	}
	filter := nostr.Filter{
		Kinds:   []int{4},
		Authors: []string{config.UserConfig.User},
		Limit:   1,
	}
	nostrSub = nostrRelay.Subscribe(nostrContext, nostr.Filters{filter})

	go func() {
		<-nostrSub.EndOfStoredEvents
		logger.Log.Debug("EOSE")
		// TODO - should I do anything with this?
	}()

	go nostrHandleMessages()
}

func nostrHandleMessages() {
	for ev := range nostrSub.Events {
		logger.Log.WithField("Event", ev).Debug("Message received")
		content, err := nip04.Decrypt(ev.Content, nostrSharedSecret)
		if err != nil {
			logger.Log.WithError(err).Error("Error decrypting contnet")
			continue
		}
		parts := strings.SplitN(content, "---", 2)
		if len(parts) != 2 {
			logger.Log.WithField("Content", content).Error("Ignoring message because it is an incorrect format")
			continue
		}
		ClipReceived(parts[1], parts[0])
	}
}

func nostrPostClip(clip string) bool {
	encrypted, err := nip04.Encrypt(fmt.Sprintf("%s---%s", config.UserConfig.Client, clip), nostrSharedSecret)
	if err != nil {
		logger.Log.WithError(err).Error("Unable to encrypt clipboard")
		return false
	}
	event := nostr.Event{
		PubKey:    config.UserConfig.User,
		CreatedAt: time.Now(),
		Kind:      4,
		Content:   encrypted,
	}
	event.Sign(config.UserConfig.Pass)
	status := nostrRelay.Publish(nostrContext, event)
	return status.String() == "success"
}

func nostrStreamClose() {
	logger.Log.Debug("Closing nostr stream")
	nostrSub.Unsub()
	nostrRelay.Close()
	nostrContextCancel()
}

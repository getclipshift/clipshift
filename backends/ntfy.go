package backends

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jhotmann/clipshift/config"
	"github.com/jhotmann/clipshift/logger"
	"github.com/sirupsen/logrus"
)

var (
	ntfyStream *resty.Response
	ntfyClient *resty.Client
)

func ntfyInit() {
	logger.Log.WithFields(logrus.Fields{
		"Host":  config.UserConfig.Host,
		"Topic": config.UserConfig.Topic,
		"User":  config.UserConfig.User,
	}).Info("Connecting to ntfy")
	ntfyClient = resty.New()
	if config.UserConfig.User != "" && config.UserConfig.Pass != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.UserConfig.User, config.UserConfig.Pass)))
		ntfyClient.SetHeader("Authorization", fmt.Sprintf("Basic %s", auth))
	}

	ntfyStreamOpen()
}

func ntfyStreamOpen() {
	var err error
	ntfyStream, err = ntfyClient.R().
		SetDoNotParseResponse(true).
		Get(fmt.Sprintf("%s/%s/json", config.UserConfig.Host, config.UserConfig.Topic))
	if err != nil {
		if strings.Contains(err.Error(), "no such host") {
			logger.Log.WithField("Error", err.Error()).Error("Waiting 5 seconds and trying again")
			time.Sleep(5 * time.Second)
			ntfyStreamOpen()
		} else {
			logger.Log.WithField("Error", err.Error()).Error("Waiting 30 seconds and trying again")
			time.Sleep(30 * time.Second)
			ntfyStreamOpen()
		}
	}
	go ntfyHandleMessages()
}

func ntfyHandleMessages() {
	defer ntfyStreamReconnect()
	scanner := bufio.NewScanner(ntfyStream.RawResponse.Body)
	for scanner.Scan() {
		var parsed NtfyMessage
		err := json.Unmarshal([]byte(scanner.Text()), &parsed)
		if err != nil {
			logger.Log.WithField("Error", err.Error()).Error("Error parsing message")
		} else if parsed.Event == NtfyEventTypes.Message {
			ClipReceived(parsed.Message, parsed.Title)
		} else {
			logger.Log.WithField("Event", parsed.Event).Debug("Response received")
		}
	}
}

func ntfyStreamReconnect() {
	ntfyStreamClose()
	ntfyStreamOpen()
}

func ntfyStreamClose() {
	logger.Log.Debug("Closing stream")
	ntfyClient.GetClient().CloseIdleConnections()
}

func ntfyPostClip(clip string) bool {
	resp, err := ntfyClient.R().
		SetHeader("Title", config.UserConfig.Client).
		SetHeader("Priority", "1").
		SetBody(clip).
		Post(fmt.Sprintf("%s/%s", config.UserConfig.Host, config.UserConfig.Topic))
	if err == nil && resp.StatusCode() == 200 {
		return true
	}
	return false
}

type NtfyMessage struct {
	Id      string `json:"id"`
	Time    int    `json:"time"`
	Event   string `json:"event"`
	Topic   string `json:"topic"`
	Message string `json:"message"`
	Title   string `json:"title"`
}

var NtfyEventTypes = struct {
	Open        string
	KeepAlive   string
	Message     string
	PollRequest string
}{
	Open:        "open",
	KeepAlive:   "keepalive",
	Message:     "message",
	PollRequest: "poll_request",
}

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
	"github.com/jhotmann/clipshift/ui"
	"github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
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
			fmt.Println(err.Error(), "\nWaiting 5 seconds and trying again")
			time.Sleep(5 * time.Second)
			ntfyStreamOpen()
		} else {
			fmt.Println(err.Error(), "\nWaiting 30 seconds and trying again")
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
			fmt.Println("Error parsing message:", err)
		} else if parsed.Event == NtfyEventTypes.Message {
			if parsed.Title == config.UserConfig.Client || parsed.Message == lastReceived {
				continue
			}
			lastReceived = parsed.Message

			if encryptionEnabled {
				lastReceived = decryptString(lastReceived)
			}

			clipboard.Write(clipboard.FmtText, []byte(lastReceived))
			fmt.Printf("Clipboard received from %s: %s\n", parsed.Title, lastReceived)
			ui.TraySetTooltip(fmt.Sprintf("%s - %s", time.Now().Format("20060102 15:04:05"), parsed.Title))
		} else {
			fmt.Println(parsed.Event, "response received")
		}
	}
}

func ntfyStreamReconnect() {
	ntfyStreamClose()
	ntfyStreamOpen()
}

func ntfyStreamClose() {
	fmt.Println("Closing stream")
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

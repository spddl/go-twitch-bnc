package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"os"
	"os/signal"

	"github.com/spddl/go-twitch-ws"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile) // https://ispycode.com/GO/Logging/Setting-output-flags

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	settings := readSettings()
	SearchByte := [][]byte{}
	for _, search := range settings.Search {
		SearchByte = append(SearchByte, []byte(search))
	}

	bot, err := twitch.NewClient(twitch.Client{
		Server:      settings.TwitchServer,
		User:        settings.User,
		Oauth:       settings.Oauth,
		Debug:       settings.Debug,
		BotVerified: settings.BotVerified,
	})
	if err != nil {
		panic(err)
	}
	defer bot.Close()

	// bot.OnNoticeMessage = func(ircMsg twitch.IRCMessage) {
	// 	log.Printf("%s\n", ircMsg)
	// }

	backlog := MutexBacklog{
		msgBacklog: make(map[string][]byte),
	}

	bot.OnPrivateMessage = func(ircMsg twitch.IRCMessage) {
		message := ircMsg.Params[1]
		channel := string(ircMsg.Params[0][1:])
		// log.Printf("%s, %s: %s\n", ircMsg.Params[0][1:], ircMsg.Tags["display-name"], message)

		timestamp := append([]byte(time.Now().Format(settings.TimestampFormat)), []byte{59, 32}...) // []byte{59, 32} "; "
		name := append(ircMsg.Tags["display-name"], []byte{58, 32}...)
		line := append(timestamp, append(name, message...)...)

		backlog.append(channel, line)

		if backlog.count(channel) > settings.Backloglimit {
			backlog.remove(channel, settings.Backloglimit) // remove unnecessary lines
		}

		for _, searchValue := range SearchByte {
			if bytes.Contains(ircMsg.Params[1], searchValue) {
				if settings.Webhook != "" {
					sendWebhook(settings.Webhook, fmt.Sprintf("#%s: %s", channel, line[20:]))
				}

				if settings.Backloglimit != 0 && settings.LogDir != "" {
					saveLogFile(backlog.getChannel(channel), channel, settings.LogDir)
					backlog.reset(channel) // reset logs
				}
			}
		}
	}

	bot.OnConnect = func(status bool) {
		bot.Join(settings.Channel)
	}

	bot.Login()

	for { // ctrl - c
		<-interrupt
		bot.Close()
		os.Exit(0)
	}
}

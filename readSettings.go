package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Settings struct {
	TwitchServer    string   `json:"twitchServer"`
	User            string   `json:"user"`
	Oauth           string   `json:"oauth"`
	Channel         []string `json:"channel"`
	TimestampFormat string   `json:"timestampFormat"`
	Debug           bool     `json:"debug"`
	BotVerified     bool     `json:"botVerified"`
	Backloglimit    int      `json:"backloglimit"`
	Search          []string `json:"search"`
	Webhook         string   `json:"webhook"`
	LogDir          string   `json:"logDir"`
}

func readSettings() Settings {
	var byteValue []byte

	// we initialize our Settings
	var settings Settings

	// Open our jsonFile
	jsonFile, err := os.Open("settings.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Println(err)
	} else {

		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		// read our opened xmlFile as a byte array.
		byteValue, err = ioutil.ReadAll(jsonFile)
		if err != nil {
			log.Fatalln(err)
		}
		// we unmarshal our byteArray which contains our
		// jsonFile's content into 'Settings' which we defined above
		err = json.Unmarshal(byteValue, &settings)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// setting default values
	// if no values present
	if settings.TwitchServer == "" {
		settings.TwitchServer = "wss://irc-ws.chat.twitch.tv"
	}
	if settings.TimestampFormat == "" {
		settings.TimestampFormat = "02.01.2006 15:04:05"
	}
	if settings.LogDir == "" {
		settings.LogDir = "."
	}

	return settings
}

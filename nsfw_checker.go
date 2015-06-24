package main

import (
	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/koyachi/go-nude"
	"github.com/thoj/go-ircevent"
	"os"
)

// Check for NSFW Content
func CheckNSFW(ircbot *irc.Connection, url string) {
	fileName, err := utils.DownloadFile(url)
	if err != nil {
		return
	}
	defer os.Remove(fileName)
	isNude, err := nude.IsNude(fileName)
	if err != nil {
		return
	}
	if isNude {
		ircbot.Privmsgf(configuration.Config.Channel, "%v is NSFW", url)
	}
}

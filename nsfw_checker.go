package main

import (
	"github.com/koyachi/go-nude"
	"github.com/thoj/go-ircevent"
	"io"
	"net/http"
	"os"
	"strings"
)

// Download a file and write to disk
func DownloadFile(url string) (string, error) {
	splittedFileName := strings.Split(url, "/")
	fileName := splittedFileName[len(splittedFileName)-1]
	output, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	defer output.Close()
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if _, err := io.Copy(output, response.Body); err != nil {
		return "", err
	}
	return fileName, nil
}

// Check for NSFW Content
func CheckNSFW(ircbot *irc.Connection, url string) {
	fileName, err := DownloadFile(url)
	if err != nil {
		return
	}
	defer os.Remove(fileName)
	isNude, err := nude.IsNude(fileName)
	if err != nil {
		return
	}
	if isNude {
		ircbot.Privmsgf(Config.Channel, "%v is NSFW", url)
	}
}

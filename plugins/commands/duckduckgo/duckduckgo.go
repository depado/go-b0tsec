package duckduckgo

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const apiURL = "http://api.duckduckgo.com/?q=%s&format=json%s"

type message struct {
	Definition       string
	DefinitionSource string
	Heading          string
	AbstractText     string
	Abstract         string
	AbstractSource   string
	Image            string
	Type             string
	AnswerType       string
	Redirect         string
	DefinitionURL    string
	Answer           string
	AbstractURL      string
	Results          []relatedTopic
	RelatedTopics    []relatedTopic
}

type relatedTopic struct {
	Result string
	Icon   struct {
		URL    string
		Height interface{}
		Width  interface{}
	}
	FirstURL string
	Text     string
}

// Plugin is the duckduckgo plugin.
type Plugin struct {
}

// Get actually sends the data to the channel
func (p Plugin) Get(ircbot *irc.Connection, from string, to string, args []string) {
	res, err := p.Fetch(strings.Join(args, " "))
	if err != nil || res == "" {
		return
	}
	ircbot.Privmsg(configuration.Config.Channel, res)
}

// Fetch returns the abstract for the given query
func (p Plugin) Fetch(query string) (string, error) {
	var t message
	url := utils.EncodeURL(apiURL, query)
	err := utils.FetchURL(url, &t)
	if err != nil {
		return "", err
	}
	return t.Abstract, nil
}

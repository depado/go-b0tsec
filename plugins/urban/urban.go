package urban

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const apiURL = "http://api.urbandictionary.com/v0/define?term=%s"

type message struct {
	Tags       []string `json:"tags"`
	ResultType string   `json:"result_type"`
	List       []struct {
		Defid       int    `json:"defid"`
		Word        string `json:"word"`
		Author      string `json:"author"`
		Permalink   string `json:"permalink"`
		Definition  string `json:"definition"`
		Example     string `json:"example"`
		ThumbsUp    int    `json:"thumbs_up"`
		ThumbsDown  int    `json:"thumbs_down"`
		CurrentVote string `json:"current_vote"`
	} `json:"list"`
	Sounds []string `json:"sounds"`
}

// Plugin is the Urban dictionary plugin.
type Plugin struct {
}

// Get actually sends the data to the channel.
func (p Plugin) Get(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	res, err := p.Fetch(strings.Join(arguments[0], " "))
	if err != nil || res == "" {
		return
	}
	ircbot.Privmsg(configuration.Config.Channel, res)
}

// Fetch returns the abstract for the given query.
func (p Plugin) Fetch(query string) (string, error) {
	var t message
	url := utils.EncodeURL(apiURL, query)
	err := utils.FetchURL(url, &t)
	if err != nil {
		return "", err
	}
	if len(t.List) == 0 {
		return "", nil
	}
	return t.List[0].Definition, nil
}

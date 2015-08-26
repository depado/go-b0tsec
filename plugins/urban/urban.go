package urban

import (
	"errors"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const apiURL = "http://api.urbandictionary.com/v0/define?term=%s"

type message struct {
	Tags       []string   `json:"tags"`
	ResultType string     `json:"result_type"`
	List       []udResult `json:"list"`
	Sounds     []string   `json:"sounds"`
}

type udResult struct {
	Defid       int    `json:"defid"`
	Word        string `json:"word"`
	Author      string `json:"author"`
	Permalink   string `json:"permalink"`
	Definition  string `json:"definition"`
	Example     string `json:"example"`
	ThumbsUp    int    `json:"thumbs_up"`
	ThumbsDown  int    `json:"thumbs_down"`
	CurrentVote string `json:"current_vote"`
}

// Plugin is the Urban dictionary plugin.
type Plugin struct {
	Last          message
	CurrentResult udResult
	Current       int
}

// Get actually sends the data to the channel.
func (p *Plugin) Get(ib *irc.Connection, nick string, general bool, arguments ...[]string) {
	var err error
	if len(arguments[0]) == 1 {
		switch arguments[0][0] {
		case "moar":
			err = p.more()
			if err != nil {
				ib.Privmsg(configuration.Config.Channel, "Nothing else here.")
			} else {
				for _, mess := range utils.SplitMessage(p.CurrentResult.Definition) {
					ib.Privmsg(configuration.Config.Channel, mess)
				}
			}
			return
		case "quote":
			for _, mess := range utils.SplitMessage(p.CurrentResult.Example) {
				ib.Privmsg(configuration.Config.Channel, mess)
			}
			return
		}
	}
	err = p.fetch(strings.Join(arguments[0], " "))
	if err != nil {
		return
	}
	ib.Privmsg(configuration.Config.Channel, p.CurrentResult.Definition)
}

func (p *Plugin) sanitize() {
	p.CurrentResult.Definition = strings.Replace(p.CurrentResult.Definition, "\r\n", " ", -1)
	p.CurrentResult.Example = strings.Replace(p.CurrentResult.Example, "\r\n", " ", -1)
}

func (p *Plugin) more() error {
	if len(p.Last.List) < p.Current+2 {
		return errors.New("Nothing else here.")
	}
	p.Current++
	p.CurrentResult = p.Last.List[p.Current]
	p.sanitize()
	return nil
}

// fetch queries the API with the given query. It then fills the Plugin's struct.
func (p *Plugin) fetch(query string) error {
	var t message
	url := utils.EncodeURL(apiURL, query)
	err := utils.FetchURL(url, &t)
	if err != nil {
		return err
	}
	if len(t.List) == 0 {
		return errors.New("No result")
	}
	p.Last = t
	p.Current = 0
	p.CurrentResult = t.List[p.Current]
	p.sanitize()
	return nil
}

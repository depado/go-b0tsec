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

// Help must send the help about this plugin.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    Allows to search a term on the Urban Dictionnary")
	ib.Privmsg(from, "    Optional argument : moar - Allows to search another definition from the previous search.")
	ib.Privmsg(from, "    Optional argument : quote - Allows to search a quote from the previous search.")
}

// Get actually sends the data to the channel.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		var err error
		if len(args) == 1 {
			switch args[0] {
			case "moar":
				err = p.more()
				if err != nil {
					utils.Send(ib, err.Error())
				} else {
					for _, m := range utils.SplitMessage(p.CurrentResult.Definition) {
						utils.Send(ib, m)
					}
				}
				return
			case "quote":
				for _, m := range utils.SplitMessage(p.CurrentResult.Example) {
					utils.Send(ib, m)
				}
				return
			}
		}
		err = p.fetch(strings.Join(args, " "))
		if err != nil {
			return
		}
		ib.Privmsg(configuration.Config.Channel, p.CurrentResult.Definition)
	}
}

// sanitize removes the \r\n escaping chars from the definitions and example of
// p.CurrentResult
func (p *Plugin) sanitize() {
	p.CurrentResult.Definition = strings.Replace(p.CurrentResult.Definition, "\r\n", " ", -1)
	p.CurrentResult.Example = strings.Replace(p.CurrentResult.Example, "\r\n", " ", -1)
}

// more switches the p.CurrentResult to the next entry or returns an error if it
// was already the last entry available.
func (p *Plugin) more() error {
	if len(p.Last.List) < p.Current+2 {
		return errors.New("Nothing else here.")
	}
	p.Current++
	p.CurrentResult = p.Last.List[p.Current]
	p.sanitize()
	return nil
}

// fetch queries the apiURL with the given query and populates the p Plugin.
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

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}

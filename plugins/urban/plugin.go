package urban

import (
	"errors"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "ud"
	apiURL        = "http://api.urbandictionary.com/v0/define?term=%s"
)

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
	Started       bool
	Last          message
	CurrentResult udResult
	Current       int
}

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

// Help must send the help about this plugin.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "Allows to search a term on the Urban Dictionnary")
	ib.Privmsg(from, "Optional argument : moar - Allows to search another definition from the previous search.")
	ib.Privmsg(from, "Optional argument : quote - Allows to search a quote from the previous search.")
}

// Get actually sends the data to the channel.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
		return
	}
	if len(args) > 0 {
		var err error
		if len(args) == 1 {
			switch args[0] {
			case "moar":
				err = p.more()
				if err != nil {
					ib.Privmsg(to, err.Error())
				} else {
					for _, m := range utils.SplitMessage(p.CurrentResult.Definition) {
						ib.Privmsg(to, m)
					}
				}
				return
			case "quote":
				for _, m := range utils.SplitMessage(p.CurrentResult.Example) {
					ib.Privmsg(to, m)
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

// Start starts the plugin and returns any occured error, nil otherwise
func (p *Plugin) Start() error {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		p.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occured error, nil otherwise
func (p *Plugin) Stop() error {
	p.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (p *Plugin) IsStarted() bool {
	return p.Started
}

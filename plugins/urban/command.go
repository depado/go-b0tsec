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
	command = "ud"
	apiURL  = "http://api.urbandictionary.com/v0/define?term=%s"
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

// Command is the Urban dictionary plugin.
type Command struct {
	Started       bool
	Last          message
	CurrentResult udResult
	Current       int
}

func init() {
	plugins.Commands[command] = new(Command)
}

// Help must send the help about this plugin.
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Allows to search a term on the Urban Dictionnary")
	ib.Privmsg(from, "Optional argument : moar - Allows to search another definition from the previous search.")
	ib.Privmsg(from, "Optional argument : quote - Allows to search a quote from the previous search.")
}

// Get actually sends the data to the channel.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		var err error
		if len(args) == 1 {
			switch args[0] {
			case "moar":
				err = c.more()
				if err != nil {
					ib.Privmsg(to, err.Error())
				} else {
					for _, m := range utils.SplitMessage(c.CurrentResult.Definition) {
						ib.Privmsg(to, m)
					}
				}
				return
			case "quote":
				for _, m := range utils.SplitMessage(c.CurrentResult.Example) {
					ib.Privmsg(to, m)
				}
				return
			}
		}
		err = c.fetch(strings.Join(args, " "))
		if err != nil {
			return
		}
		ib.Privmsg(configuration.Config.Channel, c.CurrentResult.Definition)
	}
}

// sanitize removes the \r\n escaping chars from the definitions and example of
// c.CurrentResult
func (c *Command) sanitize() {
	c.CurrentResult.Definition = strings.Replace(c.CurrentResult.Definition, "\r\n", " ", -1)
	c.CurrentResult.Example = strings.Replace(c.CurrentResult.Example, "\r\n", " ", -1)
}

// more switches the c.CurrentResult to the next entry or returns an error if it
// was already the last entry available.
func (c *Command) more() error {
	if len(c.Last.List) < c.Current+2 {
		return errors.New("Nothing else here.")
	}
	c.Current++
	c.CurrentResult = c.Last.List[c.Current]
	c.sanitize()
	return nil
}

// fetch queries the apiURL with the given query and populates the p Command.
func (c *Command) fetch(query string) error {
	var t message
	url := utils.EncodeURL(apiURL, query)
	err := utils.FetchURL(url, &t)
	if err != nil {
		return err
	}
	if len(t.List) == 0 {
		return errors.New("No result")
	}
	c.Last = t
	c.Current = 0
	c.CurrentResult = t.List[c.Current]
	c.sanitize()
	return nil
}

// Start starts the plugin and returns any occurred error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(command, configuration.Config.Commands) {
		c.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occurred error, nil otherwise
func (c *Command) Stop() error {
	c.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (c *Command) IsStarted() bool {
	return c.Started
}

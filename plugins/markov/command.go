package markov

import (
	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	command = "markov"
)

// Command is the markov.Command type
type Command struct {
	Started bool
}

func init() {
	plugins.Commands[command] = new(Command)
}

// Help displays the help for the plugin
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "This command generates a random sentence using the markov chains.")
	ib.Privmsg(from, "Example : !markov")
	ib.Privmsg(from, "Example with a target : !markov > nickname")
}

// Get actually acts
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		if i, ok := utils.IndexStringInSlice(">", args); ok && len(args) > i+1 {
			ib.Privmsgf(to, "%v: %v", args[i+1], MainChain.Generate())
		}
		return
	}
	ib.Privmsg(to, MainChain.Generate())
}

// Start starts the plugin and returns any occured error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(command, configuration.Config.Commands) {
		c.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occured error, nil otherwise
func (c *Command) Stop() error {
	c.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (c *Command) IsStarted() bool {
	return c.Started
}

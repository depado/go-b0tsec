package cleverbot

import (
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	command = "clever"
)

// Command is the cleverbot.Command type
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
	ib.Privmsg(from, "This command asks something to the Cleverbot API.")
	ib.Privmsg(from, "Example : !clever What time is it ?")
}

// Get actually acts
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if to == configuration.Config.BotName {
		to = from
	}
	if len(args) > 0 {
		answer, err := Clever.Query(strings.Join(args, " "))
		if err != nil {
			log.Println(err)
			return
		}
		if to == from {
			ib.Privmsg(to, answer)
		} else {
			ib.Privmsgf(to, "%s: %s", from, answer)
		}
	}
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

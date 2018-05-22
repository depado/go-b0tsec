package anon

import (
	"math/rand"
	"strings"

	"github.com/thoj/go-ircevent"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
)

const (
	command = "anon"
)

// Command is the anon plugin. Exposed as anon.Command.
type Command struct {
	Started bool
}

func init() {
	plugins.Commands[command] = new(Command)
}

// Help provides some help on the usage of the plugin.
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Allows to send anonymous messages on the channel where the bot is connected.")
	ib.Privmsgf(from, "Example : /msg %s !anon Hello everyone.", configuration.Config.BotName)
}

// Get actually sends the data to the channel.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		ib.Privmsgf(configuration.Config.Channel, "[%s] %v", string(from[rand.Intn(len(from))]), strings.Join(args, " "))
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

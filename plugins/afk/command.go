package afk

import (
	"strings"
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	command = "afk"
)

// Command is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Command struct {
	Started bool
}

func init() {
	plugins.Commands[command] = new(Command)
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Tell the world you're afk, for a reason. Or not.")
	ib.Privmsg(from, "Example : !afk reason.")
}

// Get is the actual call to your plugin.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	reason := ""
	if len(args) > 0 {
		reason = strings.Join(args, " ")
	}
	Map[from] = Data{time.Now(), reason}
	if reason != "" {
		ib.Privmsgf(configuration.Config.Channel, "%v is afk : %v", from, reason)
	} else {
		ib.Privmsgf(configuration.Config.Channel, "%v is afk.", from)
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

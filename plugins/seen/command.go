package seen

import (
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/plugins/afk"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	command = "seen"
)

// Command is the plugin struct. It will be exposed as packagename.Command to keep the API stable and friendly.
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
	ib.Privmsg(from, "Displays if someone is afk or the time since their last message.")
	ib.Privmsg(from, "Example : !seen nickname")
}

// Get is the actual call to your plugin.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		if to == configuration.Config.BotName {
			to = from
		}
		for _, v := range args {
			if d, ok := afk.Map[v]; ok {
				if d.Reason != "" {
					ib.Privmsgf(to, "%v has been afk for %v : %v", v, time.Since(d.Since).String(), d.Reason)
				} else {
					ib.Privmsgf(to, "%v has been afk for %v", v, time.Since(d.Since).String())
				}
			} else {
				if d, ok := Map[v]; ok {
					ib.Privmsgf(to, "Last message from %v : %v ago.", v, time.Since(d).String())
				}
			}
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

package choice

import (
	"math/rand"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "choice"
)

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct {
	Started bool
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "Makes a choice for you. Example : !choice rhum whisky vodka beer")
}

// Get is the actual call to your plugin.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
		return
	}
	if to == configuration.Config.BotName {
		to = from
	}
	if len(args) > 1 {
		c := rand.Intn(len(args))
		if to == from {
			ib.Privmsgf(to, "%v", args[c])
		} else {
			ib.Privmsgf(to, "%v: %v", from, args[c])
		}
	}
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

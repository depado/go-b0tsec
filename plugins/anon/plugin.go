package anon

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "anon"
)

// Plugin is the anon plugin. Exposed as anon.Plugin.
type Plugin struct {
	Started bool
}

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

// Help provides some help on the usage of the plugin.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "Allows to send anonymous messages on the channel where the bot is connected.")
	ib.Privmsgf(from, "Example : /msg %s !anon Hello everyone.", configuration.Config.BotName)
}

// Get actually sends the data to the channel.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
		return
	}
	if len(args) > 0 {
		ib.Privmsgf(configuration.Config.Channel, "[A] %v", strings.Join(args, " "))
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

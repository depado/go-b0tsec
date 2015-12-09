package markov

import (
	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "markov"
)

// Plugin is the markov.Plugin type
type Plugin struct {
	Started bool
}

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

// Help displays the help for the plugin
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "This command generates a random sentence using the markov chains.")
	ib.Privmsg(from, "Example : !markov")
	ib.Privmsg(from, "Example with a target : !markov > nickname")
}

// Get actually acts
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
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

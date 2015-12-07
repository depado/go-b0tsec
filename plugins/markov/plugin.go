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
type Plugin struct{}

func init() {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		plugins.Plugins[pluginCommand] = new(Plugin)
	}
}

// Help displays the help for the plugin
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "This command generates a random sentence using the markov chains.")
	ib.Privmsg(from, "Example : !markov")
	ib.Privmsg(from, "Example with a target : !markov > nickname")
}

// Get actually acts
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		if i, ok := utils.IndexStringInSlice(">", args); ok && len(args) > i+1 {
			ib.Privmsgf(to, "%v: %v", args[i+1], MainChain.Generate())
		}
		return
	}
	ib.Privmsg(to, MainChain.Generate())
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}

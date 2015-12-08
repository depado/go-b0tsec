package help

import (
	"strings"

	"github.com/depado/go-b0tsec/plugins"

	"github.com/thoj/go-ircevent"
)

// init initializes all the plugins and middlewares.
func init() {
	plugins.Plugins["help"] = new(Plugin)
}

// Plugin is the help plugin
type Plugin struct{}

// Get actually executes the command.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) == 0 {
		ib.Privmsg(from, "available commands (!help <command> to get more info) :")
		ib.Privmsg(from, strings.Join(plugins.ListPlugins(), " "))
	} else {
		if h, ok := plugins.Plugins[args[0]]; ok {
			h.Help(ib, from)
		}
	}
}

// Help shows the help for the plugin.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "so you need help using the help command to get help about other commands ?")
	ib.Privmsg(from, "with argument : show the help for a specific argument (e.g : !help ud).")
	ib.Privmsg(from, "without argument : shows the available commands.")
}

// Start returns nil since it is a core plugin
func (p *Plugin) Start() error {
	return nil
}

// Stop returns nil since it is a core plugin
func (p *Plugin) Stop() error {
	return nil
}

// IsStarted returns always true since it is a core plugin
func (p *Plugin) IsStarted() bool {
	return true
}

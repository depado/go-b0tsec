package help

import (
	"strings"

	"github.com/depado/go-b0tsec/plugins"

	"github.com/thoj/go-ircevent"
)

const (
	command = "help"
)

// init initializes all the plugins and middlewares.
func init() {
	plugins.Commands[command] = new(Command)
}

// Command is the help plugin
type Command struct{}

// Get actually executes the command.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) == 0 {
		ib.Privmsg(from, "available commands (!help <command> to get more info) :")
		ib.Privmsg(from, strings.Join(plugins.ListCommands(), " "))
	} else {
		if h, ok := plugins.Commands[args[0]]; ok {
			h.Help(ib, from)
		}
	}
}

// Help shows the help for the plugin.
func (c *Command) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "so you need help using the help command to get help about other commands ?")
	ib.Privmsg(from, "with argument : show the help for a specific argument (e.g : !help ud).")
	ib.Privmsg(from, "without argument : shows the available commands.")
}

// Start returns nil since it is a core plugin
func (c *Command) Start() error {
	return nil
}

// Stop returns nil since it is a core plugin
func (c *Command) Stop() error {
	return nil
}

// IsStarted returns always true since it is a core plugin
func (c *Command) IsStarted() bool {
	return true
}

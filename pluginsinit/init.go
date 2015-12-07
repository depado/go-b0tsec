package pluginsinit

import (
	"strings"

	"github.com/depado/go-b0tsec/plugins"
	// Blank import for init purpose by self plugin registering
	_ "github.com/depado/go-b0tsec/plugins/afk"
	_ "github.com/depado/go-b0tsec/plugins/anon"
	_ "github.com/depado/go-b0tsec/plugins/choice"
	_ "github.com/depado/go-b0tsec/plugins/define"
	_ "github.com/depado/go-b0tsec/plugins/dice"
	_ "github.com/depado/go-b0tsec/plugins/duckduckgo"
	_ "github.com/depado/go-b0tsec/plugins/github"
	_ "github.com/depado/go-b0tsec/plugins/karma"
	_ "github.com/depado/go-b0tsec/plugins/logger"
	_ "github.com/depado/go-b0tsec/plugins/markov"
	_ "github.com/depado/go-b0tsec/plugins/seen"
	_ "github.com/depado/go-b0tsec/plugins/title"
	_ "github.com/depado/go-b0tsec/plugins/translate"
	_ "github.com/depado/go-b0tsec/plugins/urban"
	_ "github.com/depado/go-b0tsec/plugins/usercommand"
	_ "github.com/depado/go-b0tsec/plugins/youtube"
	"github.com/thoj/go-ircevent"
)

// init initializes all the plugins and middlewares.
func init() {
	plugins.Plugins["help"] = new(Help)
}

// Help is the help plugin. Builtin.
type Help struct{}

// Get actually executes the command.
func (h Help) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) == 0 {
		ib.Privmsg(from, "Available commands (!help <command> to get more info) :")
		keys := []string{}
		for k := range plugins.Plugins {
			keys = append(keys, k)
		}
		ib.Privmsg(from, strings.Join(keys, " "))
	} else {
		if p, ok := plugins.Plugins[args[0]]; ok {
			p.Help(ib, from)
		}
	}
}

// Help shows the help for the plugin.
func (h Help) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "So you need help using the help command to get help about other commands ?")
	ib.Privmsg(from, "With argument : Show the help for a specific argument (e.g : !help ud).")
	ib.Privmsg(from, "Without argument : Shows the available commands.")
}

package plugins

import (
	"github.com/depado/go-b0tsec/plugins/commands/anon"
	"github.com/depado/go-b0tsec/plugins/commands/duckduckgo"
	"github.com/depado/go-b0tsec/plugins/commands/karma"
	"github.com/depado/go-b0tsec/plugins/commands/urban"
	"github.com/depado/go-b0tsec/plugins/middlewares/github"
	"github.com/depado/go-b0tsec/plugins/middlewares/logger"
	"github.com/depado/go-b0tsec/plugins/middlewares/markov"
	"github.com/thoj/go-ircevent"
)

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, string, []string)
	Help(*irc.Connection, string)
}

// Middleware represents a single middleware.
type Middleware interface {
	Get(*irc.Connection, string, string, string)
}

// Plugins is a map with the command as key and a function to be executed as value.
var Plugins = map[string]Plugin{}

// Middlewares is a slice containing the plugins that should be executed on each message reception.
var Middlewares = []func(*irc.Connection, string, string, string){}

// RegisterCommand registers a plugin in the cm CommandMap, associating the c command to the p Plugin.
func RegisterCommand(c string, p Plugin) {
	Plugins[c] = p
}

// RegisterMiddleware inserts a plugin inside the Middlewares slice.
func RegisterMiddleware(m Middleware) {
	Middlewares = append(Middlewares, m.Get)
}

// Init initializes all the plugins and middlewares.
func Init() {
	RegisterMiddleware(new(logger.Middleware))
	RegisterMiddleware(new(github.Middleware))
	RegisterMiddleware(new(markov.Middleware))
	RegisterCommand("ud", new(urban.Plugin))
	RegisterCommand("ddg", new(duckduckgo.Plugin))
	RegisterCommand("anon", new(anon.Plugin))
	RegisterCommand("markov", new(markov.Plugin))
	RegisterCommand("karma", karma.New())
	RegisterCommand("help", new(Help))
}

// Help is the help plugin. Builtin.
type Help struct{}

// Get actually executes the command.
func (h Help) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) == 0 {
		for k, p := range Plugins {
			ib.Privmsgf(from, "Command %s :", k)
			p.Help(ib, from)
		}
	} else {
		if p, ok := Plugins[args[0]]; ok {
			p.Help(ib, from)
		}
	}
}

// Help shows the help for the plugin.
func (h Help) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    With argument : Show the help for a specific argument (e.g : !help ud).")
	ib.Privmsg(from, "    Without argument : Shows the help for all the known commands.")
}

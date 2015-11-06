package plugins

import (
	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins/commands/anon"
	"github.com/depado/go-b0tsec/plugins/commands/dice"
	"github.com/depado/go-b0tsec/plugins/commands/duckduckgo"
	"github.com/depado/go-b0tsec/plugins/commands/karma"
	"github.com/depado/go-b0tsec/plugins/commands/seen"
	"github.com/depado/go-b0tsec/plugins/commands/urban"
	"github.com/depado/go-b0tsec/plugins/middlewares/github"
	"github.com/depado/go-b0tsec/plugins/middlewares/logger"
	"github.com/depado/go-b0tsec/plugins/mixins/afk"
	"github.com/depado/go-b0tsec/plugins/mixins/markov"
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
	cnf := configuration.Config
	for _, p := range cnf.Plugins {
		switch p {
		case "ud":
			RegisterCommand("ud", urban.NewPlugin())
		case "ddg":
			RegisterCommand("ddg", duckduckgo.NewPlugin())
		case "anon":
			RegisterCommand("anon", anon.NewPlugin())
		case "markov":
			RegisterCommand("markov", markov.NewPlugin())
		case "karma":
			RegisterCommand("karma", karma.NewPlugin())
		case "dice":
			RegisterCommand("dice", dice.NewPlugin())
		case "seen":
			RegisterCommand("seen", seen.NewPlugin())
		case "afk":
			RegisterCommand("afk", afk.NewPlugin())
		}
	}
	for _, m := range cnf.Middlewares {
		switch m {
		case "logger":
			RegisterMiddleware(logger.NewMiddleware())
		case "github":
			RegisterMiddleware(github.NewMiddleware())
		case "markov":
			RegisterMiddleware(markov.NewMiddleware())
		case "afk":
			RegisterMiddleware(afk.NewMiddleware())
		}
	}
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

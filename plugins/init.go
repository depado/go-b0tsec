package plugins

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins/afk"
	_ "github.com/depado/go-b0tsec/plugins/anon"
	_ "github.com/depado/go-b0tsec/plugins/choice"
	_ "github.com/depado/go-b0tsec/plugins/define"
	_ "github.com/depado/go-b0tsec/plugins/dice"
	_ "github.com/depado/go-b0tsec/plugins/duckduckgo"
	"github.com/depado/go-b0tsec/plugins/github"
	_ "github.com/depado/go-b0tsec/plugins/karma"
	"github.com/depado/go-b0tsec/plugins/logger"
	"github.com/depado/go-b0tsec/plugins/markov"
	"github.com/depado/go-b0tsec/plugins/seen"
	"github.com/depado/go-b0tsec/plugins/title"
	_ "github.com/depado/go-b0tsec/plugins/translate"
	_ "github.com/depado/go-b0tsec/plugins/urban"
	"github.com/depado/go-b0tsec/plugins/usercommand"
	"github.com/depado/go-b0tsec/plugins/youtube"
	"github.com/depado/go-b0tsec/pluginsinit"
	"github.com/thoj/go-ircevent"
)

// Middleware represents a single middleware.
type Middleware interface {
	Get(*irc.Connection, string, string, string)
}

// Middlewares is a slice containing the plugins that should be executed on each message reception.
var Middlewares = []func(*irc.Connection, string, string, string){}

// RegisterMiddleware inserts a plugin inside the Middlewares slice.
func RegisterMiddleware(m Middleware) {
	Middlewares = append(Middlewares, m.Get)
}

// init initializes all the plugins and middlewares.
func init() {
	cnf := configuration.Config
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
		case "seen":
			RegisterMiddleware(seen.NewMiddleware())
		case "youtube":
			RegisterMiddleware(youtube.NewMiddleware())
		case "title":
			RegisterMiddleware(title.NewMiddleware())
		case "usercommand":
			RegisterMiddleware(usercommand.NewMiddleware())
		}
	}
	pluginsinit.Plugins["help"] = new(Help)

	for k, m := range pluginsinit.Middlewares {
		Middlewares[k] = m
	}
}

// Help is the help plugin. Builtin.
type Help struct{}

// Get actually executes the command.
func (h Help) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) == 0 {
		ib.Privmsg(from, "Available commands (!help <command> to get more info) :")
		keys := []string{}
		for k := range pluginsinit.Plugins {
			keys = append(keys, k)
		}
		ib.Privmsg(from, strings.Join(keys, " "))
	} else {
		if p, ok := pluginsinit.Plugins[args[0]]; ok {
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

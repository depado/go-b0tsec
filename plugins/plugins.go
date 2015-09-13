package plugins

import (
	"github.com/depado/go-b0tsec/plugins/commands/anon"
	"github.com/depado/go-b0tsec/plugins/commands/duckduckgo"
	"github.com/depado/go-b0tsec/plugins/commands/urban"
	"github.com/depado/go-b0tsec/plugins/middlewares/github"
	"github.com/depado/go-b0tsec/plugins/middlewares/logger"
	"github.com/thoj/go-ircevent"
)

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, string, []string)
}

// Middleware represents a single middleware. More of a convenience struct than anything else.
type Middleware interface {
	Get(*irc.Connection, string, string, string)
}

// Plugins is a map with the command as key and a function to be executed as value.
var Plugins = map[string]func(*irc.Connection, string, string, []string){}

// Middlewares is a slice containing the plugins that should be executed on each message reception.
var Middlewares = []func(*irc.Connection, string, string, string){}

// RegisterCommand registers a plugin in the cm CommandMap, associating the c command to the p Plugin.
func RegisterCommand(c string, p Plugin) {
	Plugins[c] = p.Get
}

// RegisterMiddleware inserts a plugin inside the Middlewares slice.
func RegisterMiddleware(m Middleware) {
	Middlewares = append(Middlewares, m.Get)
}

// Init initializes all the plugins and middlewares.
func Init() {
	RegisterMiddleware(new(logger.Middleware))
	RegisterMiddleware(new(github.Middleware))
	RegisterCommand("ud", new(urban.Plugin))
	RegisterCommand("ddg", new(duckduckgo.Plugin))
	RegisterCommand("anon", new(anon.Plugin))
}

package plugins

import "github.com/thoj/go-ircevent"

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, string, []string)
	Help(*irc.Connection, string)
}

// Plugins is the map structure of all configured plugins
var Plugins = map[string]Plugin{}

// Middlewares is the slice of all configured middlewares Get() func
var Middlewares = []func(*irc.Connection, string, string, string){}

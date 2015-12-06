package pluginsinit

import "github.com/thoj/go-ircevent"

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, string, []string)
	Help(*irc.Connection, string)
}

var Plugins = map[string]Plugin{}
var Middlewares = []func(*irc.Connection, string, string, string){}

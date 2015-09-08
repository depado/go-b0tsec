package plugins

import (
	"github.com/depado/go-b0tsec/plugins/anon"
	"github.com/depado/go-b0tsec/plugins/duckduckgo"
	"github.com/depado/go-b0tsec/plugins/urban"
	"github.com/thoj/go-ircevent"
)

// PluginMap is a map with the command as key and a function to be executed as value.
var PluginMap = map[string]func(*irc.Connection, string, bool, ...[]string){}

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, bool, ...[]string)
}

// RegisterCmd registers a plugin in the cm CommandMap, associating the c command to the p Plugin.
func RegisterCmd(cm map[string]func(*irc.Connection, string, bool, ...[]string), c string, p Plugin) {
	cm[c] = p.Get
}

// Init initializes all the plugins.
func Init() {
	RegisterCmd(PluginMap, "ud", new(urban.Plugin))
	RegisterCmd(PluginMap, "ddg", new(duckduckgo.Plugin))
	RegisterCmd(PluginMap, "anon", new(anon.Plugin))
}

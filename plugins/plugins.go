package plugins

import (
	"github.com/depado/go-b0tsec/plugins/duckduckgo"
	"github.com/depado/go-b0tsec/plugins/urban"
	"github.com/thoj/go-ircevent"
)

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	Get(*irc.Connection, string, bool, ...[]string)
}

// Register registers a plugin in the cm CommandMap, associating the c command to the p Plugin.
func Register(cm map[string]func(*irc.Connection, string, bool, ...[]string), c string, p Plugin) {
	cm[c] = p.Get
}

// Init initializes all the plugins.
func Init(cm map[string]func(*irc.Connection, string, bool, ...[]string)) {
	Register(cm, "ud", new(urban.Plugin))
	Register(cm, "ddg", new(duckduckgo.Plugin))
}

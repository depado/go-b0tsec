package anon

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"
)

// Plugin is the anon plugin. Exposed as anon.Plugin.
// The anon plugin allows to send anonymous messages to the channel the bot is connected to.
// Example : !anon Hello World!
type Plugin struct {
}

// Get actually sends the data to the channel.
func (p Plugin) Get(ib *irc.Connection, nick string, general bool, arguments ...[]string) {
	if len(arguments[0]) > 0 {
		ib.Privmsg(configuration.Config.Channel, strings.Join(arguments[0], " "))
	}
}

package markov

import (
	"math/rand"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/markovchains"
	"github.com/thoj/go-ircevent"

	"strings"
)

// Middleware is the actual mmarkov.Middleware
type Middleware struct{}

// Get actually operates on the message
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !strings.HasPrefix(message, "!") {
		markovchains.MainChain.Build(message)
		markovchains.MainChain.Save()
		if rand.Intn(100) < 5 {
			ib.Privmsg(configuration.Config.Channel, markovchains.MainChain.Generate())
		}
	}
}

// Plugin is the markov.Plugin type
type Plugin struct{}

// Help displays the help for the plugin
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    This command generates a random sentence using the markov chains.")
}

// Get actually acts
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	ib.Privmsg(to, markovchains.MainChain.Generate())
}

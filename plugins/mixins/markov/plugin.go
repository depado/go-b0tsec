package markov

import "github.com/thoj/go-ircevent"

// Plugin is the markov.Plugin type
type Plugin struct{}

// Help displays the help for the plugin
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    This command generates a random sentence using the markov chains.")
}

// Get actually acts
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		if i, ok := stringInSlice(">", args); ok && len(args) > i+1 {
			ib.Privmsgf(to, "%v: %v", args[i+1], MainChain.Generate())
		}
		return
	}
	ib.Privmsg(to, MainChain.Generate())
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}

func stringInSlice(a string, list []string) (int, bool) {
	for i, b := range list {
		if b == a {
			return i, true
		}
	}
	return -1, false
}

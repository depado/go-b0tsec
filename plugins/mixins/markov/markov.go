package markov

import (
	"log"
	"math/rand"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/thoj/go-ircevent"

	"strings"
)

// Middleware is the actual mmarkov.Middleware
type Middleware struct{}

// Get actually operates on the message
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if from != configuration.Config.BotName {
		if !strings.HasPrefix(message, "!") && len(strings.Fields(message)) > 3 && !strings.Contains(message, configuration.Config.BotName) {
			message = strings.Replace(message, `"`, "", -1)
			MainChain.Build(message)
			if err := MainChain.Save(); err != nil {
				log.Println("Could not save to Bolt :", err)
			}
		}
		if !strings.HasPrefix(message, "!") {
			if rand.Intn(100) < 5 {
				ib.Privmsg(configuration.Config.Channel, MainChain.Generate())
			} else {
				if strings.Contains(message, configuration.Config.BotName) {
					ib.Privmsg(to, MainChain.Generate())
				}
			}
		}
	}
}

// NewMiddleware returns a new middleware
func NewMiddleware() *Middleware {
	MainChain = NewChain("main")
	if err := database.BotStorage.CreateBucket(bucketName); err != nil {
		log.Fatalf("While initializing Karma plugin : %s", err)
	}
	database.BotStorage.Get(bucketName, MainChain.Key, MainChain)
	return new(Middleware)
}

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

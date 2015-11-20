package markov

import (
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/thoj/go-ircevent"
)

// Middleware is the actual markov.Middleware
type Middleware struct{}

// Get actually operates on the message
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	cnf := configuration.Config
	if from != cnf.BotName {
		if !strings.HasPrefix(message, "!") && len(strings.Fields(message)) > 3 && !strings.Contains(message, cnf.BotName) {
			message = strings.Replace(message, `"`, "", -1)
			message = strings.Replace(message, `(`, "", -1)
			message = strings.Replace(message, `)`, "", -1)
			MainChain.Build(message)
			if err := MainChain.Save(); err != nil {
				log.Println("Could not save to Bolt :", err)
			}
		}
		if !strings.HasPrefix(message, "!") {
			if strings.Contains(message, cnf.BotName) {
				if to != cnf.BotName {
					ib.Privmsgf(to, "%v: %v", from, MainChain.Generate())
				} else {
					ib.Privmsg(from, MainChain.Generate())
				}
			}
		}
	}
}

// NewMiddleware returns a new middleware
func NewMiddleware() *Middleware {
	MainChain = NewChain("main")
	if err := database.BotStorage.CreateBucket(bucketName); err != nil {
		log.Fatalf("While initializing Markov middleware : %s", err)
	}
	database.BotStorage.Get(bucketName, MainChain.Key, MainChain)
	return new(Middleware)
}

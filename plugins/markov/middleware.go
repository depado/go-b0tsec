package markov

import (
	"fmt"
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "markov"
)

// Middleware is the actual markov.Middleware
type Middleware struct {
	Started bool
}

func init() {
	plugins.Middlewares[middlewareName] = new(Middleware)
}

// Get actually operates on the message
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
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

// Start starts the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		MainChain = NewChain("main")
		if err := database.BotStorage.CreateBucket(bucketName); err != nil {
			return fmt.Errorf("While initializing Markov middleware : %s", err)
		}
		database.BotStorage.Get(bucketName, MainChain.Key, MainChain)

		m.Started = true
	}
	return nil
}

// Stop stops the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Stop() error {
	MainChain = nil
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

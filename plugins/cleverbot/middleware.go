package cleverbot

import (
	"log"
	"strings"

	"github.com/thoj/go-ircevent"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
)

const (
	middlewareName = "clever"
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
	if to == cnf.BotName {
		to = from
	}
	if from != cnf.BotName {
		if !strings.HasPrefix(message, "!") && strings.Contains(message, cnf.BotName) {
			sp := strings.Fields(message)
			if len(sp) > 1 {
				answer, err := Clever.Query(strings.Join(strings.Fields(message)[1:], " "))
				if err != nil {
					log.Println(err)
					return
				}
				if to == from {
					ib.Privmsg(to, answer)
				} else {
					ib.Privmsgf(to, "%s: %s", from, answer)
				}
			}
		}
	}
}

// Start starts the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		m.Started = true
	}
	return nil
}

// Stop stops the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Stop() error {
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

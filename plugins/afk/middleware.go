package afk

import (
	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"

	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "afk"
)

// Middleware is the afk middleware
type Middleware struct {
	Started bool
}

func init() {
	plugins.Middlewares[middlewareName] = new(Middleware)
}

// Get actually sends the data
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
	if _, ok := Map[from]; ok {
		delete(Map, from)
		ib.Privmsgf(configuration.Config.Channel, "%v is back.", from)
	}
}

// Start starts the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		m.Started = true
	}
	return nil
}

// Stop returns nil when it didn’t encounter any error, the error otherwise
func (m *Middleware) Stop() error {
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

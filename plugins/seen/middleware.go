package seen

import (
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "seen"
)

// Middleware is the seen middleware
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
	Map[from] = time.Now()
}

// Start starts the middleware and returns any occured error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		m.Started = true
	}
	return nil
}

// Stop stops the middleware and returns any occured error, nil otherwise
func (m *Middleware) Stop() error {
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

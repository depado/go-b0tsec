package afk

import (
	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"
)

// Middleware is the github middleware
type Middleware struct{}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if _, ok := Map[from]; ok {
		delete(Map, from)
		ib.Privmsgf(configuration.Config.Channel, "%v is back.", from)
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}
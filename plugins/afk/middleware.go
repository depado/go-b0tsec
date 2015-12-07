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
type Middleware struct{}

func init() {
	m := plugins.Middlewares
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		plugins.Middlewares = append(m, new(Middleware).Get)
	}
}

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

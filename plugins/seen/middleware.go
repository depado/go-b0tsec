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
type Middleware struct{}

func init() {
	m := plugins.Middlewares
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		plugins.Middlewares = append(m, new(Middleware))
	}
}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	Map[from] = time.Now()
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

// Stop returns nil when it didnt encounter any error, the error otherwise
func (m Middleware) Stop() error {
	return nil
}

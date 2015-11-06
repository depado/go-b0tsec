package seen

import (
	"time"

	"github.com/thoj/go-ircevent"
)

// Middleware is the github middleware
type Middleware struct{}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	Map[from] = time.Now()
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

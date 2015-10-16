package markov

import "github.com/thoj/go-ircevent"

// Middleware is the actual mmarkov.Middleware
type Middleware struct{}

// Get actually operates on the message
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
}

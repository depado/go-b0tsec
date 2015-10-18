package logger

import (
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

var urlRegex, _ = regexp.Compile("^https?:.*$")

// Middleware is the actual logger.Middleware
type Middleware struct{}

// Get actually operates on the message
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	utils.HistoryLogger.Println(from + " : " + message)
	for _, field := range strings.Fields(message) {
		if urlRegex.MatchString(field) {
			utils.LinkLogger.Println(from + " : " + field)
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

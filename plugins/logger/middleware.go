package logger

import (
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/pluginsinit"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "logger"
)

var urlRegex, _ = regexp.Compile("^https?:.*$")

// Middleware is the actual logger.Middleware
type Middleware struct{}

func init() {
	m := pluginsinit.Middlewares
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		m = append(m, new(Middleware).Get)
	}
}

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

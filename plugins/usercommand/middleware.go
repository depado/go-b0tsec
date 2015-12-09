package usercommand

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"

	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "usercommand"
)

// Middleware is the usercommand middleware.
type Middleware struct {
	Started bool
}

func init() {
	plugins.Middlewares = append(plugins.Middlewares, new(Middleware))
}

// Get actually operates on the message
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
	cnf := configuration.Config
	if strings.HasPrefix(message, cnf.UserCommandCharacter) {
		c := Command{message[1:], ""}
		database.BotStorage.Get(bucketName, c.Name, &c)
		ib.Privmsg(to, c.Value)
	}
}

// Start starts the middleware and returns any occured error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		CreateBucket()
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

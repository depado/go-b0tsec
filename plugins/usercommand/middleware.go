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
	if !strings.HasPrefix(message, cnf.UserCommandCharacter) {
		return
	}

	splitted_msg := strings.Fields(message[1:])
	c := Command{splitted_msg[0], ""}
	database.BotStorage.Get(bucketName, c.Name, &c)

	if strings.HasPrefix(c.Value, cnf.CommandCharacter) {
		if len(c.Value) > 1 {
			splitted_command := strings.Fields(c.Value[1:])
			command := splitted_command[0]
			args := append(splitted_command[1:], splitted_msg[1:]...)
			if p, ok := plugins.Plugins[command]; ok {
				p.Get(ib, from, to, args)
				return
			}
		}
	}
	var msg string
	msg = c.Value
	if len(splitted_msg) > 1 {
		msg += " " + strings.Join(splitted_msg[1:], " ")
	}
	ib.Privmsg(to, msg)
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

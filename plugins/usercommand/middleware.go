package usercommand

import (
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"

	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "uc"
)

// Middleware is the usercommand middleware.
type Middleware struct {
	Started bool
}

func init() {
	plugins.Middlewares[middlewareName] = new(Middleware)
}

// Get actually operates on the message
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
	cnf := configuration.Config
	if len(message) == 1 || !strings.HasPrefix(message, cnf.UserCommandCharacter) {
		return
	}

	splittedMsg := strings.Fields(message[1:])
	uc := UserCommand{splittedMsg[0], ""}

	if err := database.BotStorage.Get(bucketName, uc.Name, &uc); err != nil {
		log.Println(err)
		return
	}

	if strings.HasPrefix(uc.Value, cnf.CommandCharacter) {
		if len(uc.Value) > 1 {
			splittedCommand := strings.Fields(uc.Value[1:])
			command := splittedCommand[0]
			args := append(splittedCommand[1:], splittedMsg[1:]...)
			if p, ok := plugins.Commands[command]; ok {
				p.Get(ib, from, to, args)
				return
			}
		}
	}
	var msg string
	msg = uc.Value
	if len(splittedMsg) > 1 {
		msg += " " + strings.Join(splittedMsg[1:], " ")
	}
	ib.Privmsg(to, msg)
}

// Start starts the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		CreateBucket()
		m.Started = true
	}
	return nil
}

// Stop stops the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Stop() error {
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

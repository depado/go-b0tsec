package usercommand

import (
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/thoj/go-ircevent"
)

type Middleware struct{}

// Get actually operates on the message
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	cnf := configuration.Config
	if strings.HasPrefix(message, cnf.UserCommandCharacter) {
		c := Command{message[1:], ""}
		database.BotStorage.Get(bucketName, c.Name, &c)

		/*	if strings.HasPrefix(c.Value, "!") {
			if len(c.Value) > 1 {
				splitted := strings.Fields(c.Value[1:])
				command := splitted[0]
				args := splitted[1:]
				if p, ok := plugins.Plugins[command]; ok {
					p.Get(ib, from, to, args)
				}
			}
		} else {*/
		ib.Privmsgf(to, c.Value)
		//}
	}
}

// NewMiddleware returns a new middleware
func NewMiddleware() *Middleware {
	CreateBucket()
	return new(Middleware)
}

package main

import (
	"crypto/tls"
	"log"
	"strings"

	"github.com/thoj/go-ircevent"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	_ "github.com/depado/go-b0tsec/pluginsinit"
	"github.com/depado/go-b0tsec/utils"
)

func main() {
	var err error

	cnf := configuration.Config

	// Loggers initialization
	if err = utils.InitLoggers(); err != nil {
		log.Fatalf("Something went wrong with the loggers : %v", err)
	}
	defer utils.HistoryFile.Close()
	defer utils.LinkFile.Close()

	// Defer all Close() methods
	defer database.BotStorage.Close()

	// Bot initialization
	ib := irc.IRC(cnf.BotName, cnf.BotName)
	if cnf.TLS {
		ib.UseTLS = true
		if cnf.InsecureTLS {
			ib.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
	if err = ib.Connect(cnf.Server); err != nil {
		log.Fatal(err)
	}

	// Callback on 'Connected' event
	ib.AddCallback("001", func(e *irc.Event) {
		ib.Join(cnf.Channel)
	})

	// Callback on 'Message' event
	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		from := e.Nick
		to := e.Arguments[0]
		m := e.Message()

		for _, c := range plugins.Middlewares {
			c(ib, from, to, m)
		}

		if strings.HasPrefix(m, cnf.CommandCharacter) {
			if len(m) > 1 {
				splitted := strings.Fields(m[1:])
				command := splitted[0]
				args := splitted[1:]
				if p, ok := plugins.Plugins[command]; ok {
					p.Get(ib, from, to, args)
				}
			}
		}
	})
	ib.Loop()
}

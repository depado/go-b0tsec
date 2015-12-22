package main

import (
	"crypto/tls"
	"flag"
	"log"
	"strings"

	"github.com/thoj/go-ircevent"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	_ "github.com/depado/go-b0tsec/pluginsinit"
)

func main() {
	var err error

	cnf := configuration.Config

	// Argument parsing
	configuration.ConfPath = flag.String("c", "conf.yml", "Local path to configuration file.")
	flag.Parse()

	// Load the configuration of the bot
	configuration.Load()

	// Storage initialization
	if err := database.BotStorage.Open(); err != nil {
		log.Fatalf("something went wrong with the databse : %v", err)
	}

	plugins.Start()

	// Defer all Close/Stop methods
	defer database.BotStorage.Close()
	defer plugins.Stop()

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
			c.Get(ib, from, to, m)
		}

		if strings.HasPrefix(m, cnf.CommandCharacter) {
			if len(m) > 1 {
				splitted := strings.Fields(m[1:])
				command := splitted[0]
				args := splitted[1:]
				if p, ok := plugins.Commands[command]; ok {
					p.Get(ib, from, to, args)
				}
			}
		}
	})
	ib.Loop()
}

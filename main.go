package main

import (
	"crypto/tls"
	"flag"
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/contentmanager"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/markovchains"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

func main() {
	var err error

	// Argument parsing
	confPath := flag.String("c", "conf/conf.yml", "Local path to configuration file.")
	noExt := flag.Bool("no-external", false, "Disable the external resource collection.")
	flag.Parse()

	// Load the configuration of the bot
	configuration.LoadConfiguration(*confPath)

	// External resource initialization if needed
	if !*noExt {
		if err = contentmanager.LoadAndStartExternalResources(); err != nil {
			log.Fatalf("Error Starting External Resources : %v", err)
		}
	}

	// Loggers initialization
	if err = utils.InitLoggers(); err != nil {
		log.Fatalf("Something went wrong with the loggers : %v", err)
	}
	defer utils.HistoryFile.Close()
	defer utils.LinkFile.Close()

	// Storage initialization
	if err = database.BotStorage.Open(); err != nil {
		log.Fatalf("Something went wrong with the databse : %v", err)
	}

	// Markov chain initialization
	if err = markovchains.Init(); err != nil {
		log.Fatalf("Something went wrong with the markov chains : %v", err)
	}

	// Bot initialization
	ib := irc.IRC(configuration.Config.BotName, configuration.Config.BotName)
	if configuration.Config.TLS {
		ib.UseTLS = true
		if configuration.Config.InsecureTLS {
			ib.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		}
	}
	ib.Connect(configuration.Config.Server)

	// Plugins initialization
	plugins.Init()

	// Callback on 'Connected' event
	ib.AddCallback("001", func(e *irc.Event) {
		ib.Join(configuration.Config.Channel)
	})

	// Callback on 'Message' event
	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		from := e.Nick
		to := e.Arguments[0]
		m := e.Message()

		for _, c := range plugins.Middlewares {
			c(ib, from, to, m)
		}

		if strings.HasPrefix(m, "!") {
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

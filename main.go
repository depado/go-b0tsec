package main

import (
	"crypto/tls"
	"flag"
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/contentmanager"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

func main() {
	// Argument parsing
	confPath := flag.String("c", "conf/conf.yml", "Local path to configuration file.")
	noExt := flag.Bool("no-external", false, "Disable the external resource collection.")
	flag.Parse()

	// Load the configuration of the bot
	configuration.LoadConfiguration(*confPath)

	// Start external resource collection if needed
	if !*noExt {
		if err := contentmanager.LoadAndStartExternalResources(); err != nil {
			log.Println("Error Starting External Resources : ", err)
		}
	}

	// Loggers Initialization
	err := utils.InitLoggers()
	if err != nil {
		log.Fatalf("Something went wrong with the loggers %v", err)
	}
	defer utils.HistoryFile.Close()
	defer utils.LinkFile.Close()

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

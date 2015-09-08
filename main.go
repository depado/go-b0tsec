package main

import (
	"crypto/tls"
	"flag"
	"log"
	"regexp"
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

	// Url Regexp
	urlRegex, _ := regexp.Compile("^https?:.*$")

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
	plugins.Init(CommandMapping)

	// Callback on 'Connected' event
	ib.AddCallback("001", func(e *irc.Event) {
		ib.Join(configuration.Config.Channel)
	})

	// Callback on 'Message' event
	ib.AddCallback("PRIVMSG", func(e *irc.Event) {
		nick := e.Nick
		message := e.Message()
		sentTo := e.Arguments[0]

		utils.HistoryLogger.Println(e.Nick + " : " + message)

		// Logging capability
		go func(message string) {
			for _, field := range strings.Fields(message) {
				if urlRegex.MatchString(field) {
					utils.LinkLogger.Println(e.Nick + " : " + field)
				}
			}
		}(message)

		if strings.HasPrefix(message, "!") {
			if len(message) > 1 {
				commandArray := strings.Fields(message[1:])
				command := commandArray[0]
				if commandCallback, ok := CommandMapping[command]; ok {
					if len(commandArray) > 1 {
						commandCallback(ib, nick, sentTo == configuration.Config.Channel, commandArray[1:])
					} else {
						commandCallback(ib, nick, sentTo == configuration.Config.Channel)
					}
				} else if response, ok := BasicsWithNickname[command]; ok {
					BasicCommandFormat(ib, nick, response)
				} else if generic, ok := GenericCommandMapping[command]; ok {
					GenericCommandFormat(ib, nick, sentTo == configuration.Config.Channel, generic, commandArray[1:])
				}
			}
		}
	})
	ib.Loop()
}

package main

import (
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
	//imgRegex, _ := regexp.Compile("^https?:.*(jpg|png|gif)$")
	urlRegex, _ := regexp.Compile("^https?:.*$")

	// Load the configuration of the bot
	// IDEA : Give the configuration as a parameter to the program.
	configuration.LoadConfiguration()
	if err := contentmanager.LoadAndStartExternalResources(); err != nil {
		log.Println("Error Starting External Resources : ", err)
	}

	// Loggers Initialization
	err := utils.InitLoggers()
	if err != nil {
		log.Fatalf("Something went wrong with the loggers %v", err)
	}
	defer utils.HistoryFile.Close()
	defer utils.LinkFile.Close()

	// Bot initialization
	ircbot := irc.IRC(configuration.Config.BotName, configuration.Config.BotName)
	ircbot.Connect(configuration.Config.Server)

	// Plugins initialization
	plugins.Init(CommandMapping)

	// Callback on 'Connected' event
	ircbot.AddCallback("001", func(e *irc.Event) {
		ircbot.Join(configuration.Config.Channel)
	})

	// Callback on 'Message' event
	ircbot.AddCallback("PRIVMSG", func(e *irc.Event) {
		nick := e.Nick
		message := e.Message()
		sentTo := e.Arguments[0]

		utils.HistoryLogger.Println(e.Nick + " : " + message)

		go func(message string) {
			for _, field := range strings.Fields(message) {
				if urlRegex.MatchString(field) {
					utils.LinkLogger.Println(e.Nick + " : " + field)
				}
				// if imgRegex.MatchString(field) {
				// 	utils.CheckNSFW(ircbot, field)
				// }
			}
		}(message)

		if strings.HasPrefix(message, "!") {
			if len(message) > 1 {
				commandArray := strings.Fields(message[1:])
				command := commandArray[0]
				if commandCallback, ok := CommandMapping[command]; ok {
					if len(commandArray) > 1 {
						commandCallback(ircbot, nick, sentTo == configuration.Config.Channel, commandArray[1:])
					} else {
						commandCallback(ircbot, nick, sentTo == configuration.Config.Channel)
					}
				} else if response, ok := BasicsWithNickname[command]; ok {
					BasicCommandFormat(ircbot, nick, response)
				} else if generic, ok := GenericCommandMapping[command]; ok {
					GenericCommandFormat(ircbot, nick, sentTo == configuration.Config.Channel, generic, commandArray[1:])
				}
			}
		}
	})
	ircbot.Loop()
}

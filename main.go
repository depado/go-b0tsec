package main

import (
	"github.com/thoj/go-ircevent"
	"regexp"
	"strings"
)

func main() {
	r, _ := regexp.Compile("^https?:.*(jpg|png|gif)$")

	LoadConfiguration()

	ircbot := irc.IRC(Config.BotName, Config.BotName)
	ircbot.Connect(Config.Server)

	ircbot.AddCallback("001", func(e *irc.Event) {
		ircbot.Join(Config.Channel)
	})

	ircbot.AddCallback("PRIVMSG", func(e *irc.Event) {
		nick := e.Nick
		message := e.Message()
		sent_to := e.Arguments[0]

		for _, field := range strings.Fields(message) {
			if r.MatchString(field) {
				go CheckNSFW(ircbot, field)
			}
		}

		if strings.HasPrefix(message, "!") {
			if len(message) > 1 {
				commandArray := strings.Fields(message[1:])
				command := commandArray[0]
				if commandCallback, ok := CommandMapping[command]; ok {
					if len(commandArray) > 1 {
						commandCallback(ircbot, nick, sent_to == Config.Channel, commandArray[1:])
					} else {
						commandCallback(ircbot, nick, sent_to == Config.Channel)
					}
				} else if response, ok := BasicsWithNickname[command]; ok {
					BasicCommandFormat(ircbot, nick, response)
				} else if generic, ok := GenericCommandMapping[command]; ok {
					GenericCommandFormat(ircbot, nick, sent_to == Config.Channel, generic, commandArray[1:])
				}
			}
		}
	})
	ircbot.Loop()
}

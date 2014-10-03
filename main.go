package main

import (
	"github.com/thoj/go-ircevent"
	"regexp"
	"strings"
)

const (
	server  = "irc.freenode.net:6667"
	channel = "#n0sec"
	botname = "gob0tsec"
)

func main() {
	r, _ := regexp.Compile("^https?:.*(jpg|png|gif)$")
	ircbot := irc.IRC(botname, botname)
	ircbot.Connect(server)

	ircbot.AddCallback("001", func(e *irc.Event) {
		ircbot.Join(channel)
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
			commandArray := strings.Fields(message[1:])
			command := commandArray[0]
			if commandCallback, ok := CommandMapping[command]; ok {
				if len(commandArray) > 1 {
					commandCallback(ircbot, nick, channel, sent_to == channel, commandArray[1:])
				} else {
					commandCallback(ircbot, nick, channel, sent_to == channel)
				}
			}
		}
	})
	ircbot.Loop()
}

package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

const (
	server  = "irc.freenode.net:6667"
	channel = "#n0sec"
	botname = "gob0tsec"
)

func main() {
	ircbot := irc.IRC(botname, botname)
	ircbot.Connect(server)

	ircbot.AddCallback("001", func(e *irc.Event) {
		ircbot.Join(channel)
	})

	ircbot.AddCallback("PRIVMSG", func(e *irc.Event) {
		nick := e.Nick
		message := e.Message()
		sent_to := e.Arguments[0]
		if strings.HasPrefix(message, "!") {
			commandArray := strings.Fields(message[1:])
			command := commandArray[0]
			if commandCallback, ok := commandMapping[command]; ok {
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

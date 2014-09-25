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

func handle_chan_command(ircbot *irc.Connection, command, nick string, arguments []string) {

}

func handle_private_command(ircbot *irc.Connection, command, nick string, arguments []string) {

}

func handle_private_message(ircbot *irc.Connection, message, nick string) {
	ircbot.Privmsg(channel, message)
}

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
		if strings.HasPrefix(string(message[0]), "!") {
			arguments := strings.Split(message, " ")
			if sent_to == botname {
				handle_private_command(ircbot, message, nick, arguments)
			} else if sent_to == channel {
				handle_chan_command(ircbot, message, nick, arguments)
			}
		} else if sent_to == botname {
			handle_private_message(ircbot, message, nick)
		}
	})
	ircbot.Loop()
}

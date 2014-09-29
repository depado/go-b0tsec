package main

import (
	"github.com/thoj/go-ircevent"
)

var commandMapping = map[string]func(*irc.Connection, string, string, bool, ...[]string){
	"command": func(ircbot *irc.Connection, nick, channel string, general bool, arguments ...[]string) {
		if len(arguments) > 0 {
			ircbot.Privmsgf(channel, "Got 'command' from %v with arguments %v", nick, arguments[0])
		} else {
			ircbot.Privmsgf(channel, "Got 'command' from %v without arguments", nick)
		}
	},
}

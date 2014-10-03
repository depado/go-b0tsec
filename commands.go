package main

import (
	"github.com/thoj/go-ircevent"
	"strings"
)

/*
   Arguments in the functions below may be an empty array or an array containing one array of strings.
   This is used only to state that the string array can be optionnal (some commands won't take any arguments)
   For commands who don't take arguments, just ignore them.
*/

var CommandMapping = map[string]func(*irc.Connection, string, string, bool, ...[]string){
	"afk": func(ircbot *irc.Connection, nick, channel string, general bool, arguments ...[]string) {
		if len(arguments) > 0 {
			ircbot.Privmsgf(channel, "%v is afk : %v", nick, strings.Join(arguments[0], " "))
		} else {
			ircbot.Privmsgf(channel, "%v is afk.", nick)
		}
	},
	"strapon": strapon,
	"strpn":   strapon,
}

func strapon(ircbot *irc.Connection, nick, channel string, general bool, arguments ...[]string) {
	state := "It succeed"
	if len(arguments) > 0 {
		ircbot.Privmsgf(channel, "%v uses a strapon on %v. %v.", nick, strings.Join(arguments[0], ", "), state)
	} else {
		ircbot.Privmsgf(channel, "%v uses a strapon on the whole room. %v.", nick, state)
	}
}

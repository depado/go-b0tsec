package main

import (
	"github.com/thoj/go-ircevent"
	"math/rand"
	"strings"
	"time"
)

/*
   Arguments in the functions below may be an empty array or an array containing one array of strings.
   This is used only to state that the string array can be optionnal (some commands won't take any arguments)
   For commands who don't take arguments, just ignore them.
*/
var CommandMapping = map[string]func(*irc.Connection, string, bool, ...[]string){
	"afk":       afk,
	"strapon":   strapon,
	"strpn":     strapon,
	"eightball": eightball,
	"bifle":     bifle,
}

func bifle(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	if len(arguments) > 0 {
		ircbot.Privmsgf(Config.Channel, "%v slaps %v with his cock !", nick, GenerateTargetString(arguments[0]))
	} else {
		ircbot.Privmsgf(Config.Channel, "%v slaps the whole room with his cock.", nick)
	}
}

func afk(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	if len(arguments) > 0 {
		ircbot.Privmsgf(Config.Channel, "%v is afk : %v", nick, strings.Join(arguments[0], " "))
	} else {
		ircbot.Privmsgf(Config.Channel, "%v is afk.", nick)
	}
}

func strapon(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	state := "It succeed."
	if len(arguments) > 0 {
		ircbot.Privmsgf(Config.Channel, "%v uses a strapon on %v. %v.", nick, GenerateTargetString(arguments[0]), state)
	} else {
		ircbot.Privmsgf(Config.Channel, "%v uses a strapon on the whole room. %v.", nick, state)
	}
}

func eightball(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	rand.Seed(time.Now().UnixNano())
	ircbot.Privmsg(Config.Channel, EightBallAnswers[rand.Intn(len(EightBallAnswers))])
}

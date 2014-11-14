package main

import (
	"github.com/thoj/go-ircevent"
	"math/rand"
	"strings"
	"time"
)

type GenericCommand struct {
	WithTargets string
	NoTargets   string
}

/*
   Arguments in the functions below may be an empty array or an array containing one array of strings.
   This is used only to state that the string array can be optionnal (some commands won't take any arguments)
   For commands who don't take arguments, just ignore them.
*/
var CommandMapping = map[string]func(*irc.Connection, string, bool, ...[]string){
	"afk":       afk,
	"eightball": eightball,
	"say":       say,
}

func afk(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	if len(arguments) > 0 {
		ircbot.Privmsgf(Config.Channel, "%v is afk : %v", nick, strings.Join(arguments[0], " "))
	} else {
		ircbot.Privmsgf(Config.Channel, "%v is afk.", nick)
	}
}

func eightball(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	rand.Seed(time.Now().UnixNano())
	ircbot.Privmsg(Config.Channel, EightBallAnswers[rand.Intn(len(EightBallAnswers))])
}

func say(ircbot *irc.Connection, nick string, general bool, arguments ...[]string) {
	if len(arguments[0]) > 0 {
		ircbot.Privmsg(Config.Channel, strings.Join(arguments[0], " "))
	} else {
		ircbot.Privmsg(nick, "This function needs a string.")
	}
}

/*
	This map associates a string (command) with a struct which defines formatted messages to be sent.
	This map is associated with the GenericCommandFormat function
	Possible enhancement : Adding different behaviour in the struct.
*/
var GenericCommandMapping = map[string]GenericCommand{
	"bj": GenericCommand{
		WithTargets: "%v gives %v a blowjob.",
		NoTargets:   "%v tries to suck his own cock but fails.",
	},
	"bifle": GenericCommand{
		WithTargets: "%v slaps %v with his cock !",
		NoTargets:   "%v slaps the whole room with his cock.",
	},
	"strapon": GenericCommand{
		WithTargets: "%v uses a strapon on %v.",
		NoTargets:   "%v uses a strapon on the whole room.",
	},
}

func GenericCommandFormat(ircbot *irc.Connection, nick string, general bool, generic GenericCommand, arguments ...[]string) {
	if len(arguments) > 0 {
		ircbot.Privmsgf(Config.Channel, generic.WithTargets, nick, GenerateTargetString(arguments[0]))
	} else {
		ircbot.Privmsgf(Config.Channel, generic.NoTargets, nick)
	}
}

/*
	This commands are designed to only display a message. No arguments needs to be given.
	(Only the nickname is known and used as the string formatter)
	This map is associated with the BasicCommandFormat function.
*/
var BasicsWithNickname = map[string]string{
	"nom":    "is going to eat.",
	"smoke":  "is going to smoke.",
	"drug":   "is going to smoke a big joint.",
	"coffee": "is going to drink a coffee.",
}

func BasicCommandFormat(ircbot *irc.Connection, nick, response string) {
	ircbot.Privmsgf(Config.Channel, "%v %v", nick, response)
}

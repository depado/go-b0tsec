package usercommand

import (
	"log"
	"strings"

	"github.com/thoj/go-ircevent"
)

// Plugin is the usercommand.Plugin type
type Plugin struct{}

// Help displays the help for the plugin
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "This command allows user to create, list and delete their own commands")
	ib.Privmsg(from, "Example :")
	ib.Privmsg(from, "!uc introduce Hi im go-b0tsec !")
	ib.Privmsg(from, "    Command introduce added")
	ib.Privmsg(from, ".introduce")
	ib.Privmsg(from, "    Hi im go-b0tsec !")
	ib.Privmsg(from, "!uc")
	ib.Privmsg(from, "    introduce")
	ib.Privmsg(from, "!uc introduce")
	ib.Privmsg(from, "    Command introduce deleted")
}

// Get actually acts
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 1 {
		// Setting a command
		c := Command{Name: args[0], Value: strings.Join(args[1:], " ")}
		if err := c.Save(); err != nil {
			log.Println("Could not save to Bolt : ", err)
			return
		}
		ib.Privmsgf(to, "Command %s added", c.Name)
		return
	}
	if len(args) == 1 {
		// Removes the command
		c := Command{Name: args[0]}
		if err := c.Delete(); err != nil {
			log.Println("Could not delete Bolt data : ", err)
			return
		}
		ib.Privmsgf(to, "Command %s deleted", c.Name)
		return
	}
	// List saved commands
	l := make([]string, 0)
	if err := List(&l); err != nil {
		log.Println("Error during listing : ", err)
	}
	if len(l) < 1 {
		ib.Privmsg(to, "No user command registered in db.")
		return
	}
	ib.Privmsgf(to, "Registered commands : %s", strings.Join(l, " "))
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	CreateBucket()
	return new(Plugin)
}

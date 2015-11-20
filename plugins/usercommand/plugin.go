package usercommand

import (
	"fmt"
	"log"
	"strings"

	"github.com/thoj/go-ircevent"
)

// Plugin is the usercommand.Plugin type
type Plugin struct{}

// Help displays the help for the plugin
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "This command allows user to bind their own command to :")
	ib.Privmsg(from, " - text printing")
	//	ib.Privmsg(from, " - command alias")
	ib.Privmsg(from, "Example : !set introduce Hi im go-b0tsec !")
	ib.Privmsg(from, "          .introduce will call the bot to print the text above.")
	//	ib.Privmsg(from, "Example : !set 421 !dice 3d6")
	//	ib.Privmsg(from, "          .421 will be interpreted as if the command !dice 3d6 was called")
}

// Get actually acts
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 1 {
		// Setting a command
		c := Command{Name: args[0], Value: strings.Join(args[1:], " ")}
		if err := c.Save(); err != nil {
			log.Println("Could not save to Bolt :", err)
			return
		}
		m := fmt.Sprintf("Command %s added", c.Name)
		ib.Privmsg(to, m)
	}
	if len(args) == 1 {
		// Removes the command
	}
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	CreateBucket()
	return new(Plugin)
}

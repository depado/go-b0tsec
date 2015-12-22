package usercommand

import (
	"log"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	command = "uc"
)

func init() {
	plugins.Commands[command] = new(Command)
}

// Command is the usercommand.Command type
type Command struct {
	Started bool
}

// Help displays the help for the plugin
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "This command allows user to create, list and delete their own commands")
	ib.Privmsg(from, "Example :")
	ib.Privmsg(from, "!uc introduce Hi im go-b0tsec !")
	ib.Privmsg(from, "> Command introduce added")
	ib.Privmsg(from, ".introduce")
	ib.Privmsg(from, "> Hi im go-b0tsec !")
	ib.Privmsg(from, "!uc")
	ib.Privmsg(from, "> introduce")
	ib.Privmsg(from, "!uc introduce")
	ib.Privmsg(from, "> Command introduce deleted")
}

// Get actually acts
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if to == configuration.Config.BotName {
		to = from
	}
	if len(args) > 1 {
		// Setting a command
		uc := UserCommand{Name: args[0], Value: strings.Join(args[1:], " ")}
		if err := uc.Save(); err != nil {
			log.Println("Could not save to Bolt : ", err)
			return
		}
		ib.Privmsgf(to, "Command %s added", uc.Name)
		return
	}
	if len(args) == 1 {
		// Removes the command
		uc := UserCommand{Name: args[0]}
		if err := uc.Delete(); err != nil {
			log.Println("Could not delete Bolt data : ", err)
			return
		}
		ib.Privmsgf(to, "Command %s deleted", uc.Name)
		return
	}
	// List saved commands
	var l []string
	if err := List(&l); err != nil {
		log.Println("Error during listing : ", err)
	}
	if len(l) < 1 {
		ib.Privmsg(to, "No user command registered.")
		return
	}
	ib.Privmsgf(to, "Registered commands : %s", strings.Join(l, " "))
}

// Start starts the plugin and returns any occured error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(command, configuration.Config.Commands) {
		CreateBucket()
		c.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occured error, nil otherwise
func (c *Command) Stop() error {
	c.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (c *Command) IsStarted() bool {
	return c.Started
}

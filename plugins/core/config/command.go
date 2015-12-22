package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"

	"github.com/thoj/go-ircevent"
)

const (
	command = "config"
)

type modifier struct {
	Commands    []string
	Middlewares []string
}

// Command is the help plugin
type Command struct {
	pending bool
	auth    bool
	toStart *modifier
	toStop  *modifier
	args    []string
}

// init initializes all the plugins and middlewares.
func init() {
	plugins.Commands[command] = new(Command)
}

// Get actually executes the command.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if to == configuration.Config.BotName {
		to = from
	}
	if !utils.StringInSlice(from, configuration.Config.Admins) {
		ib.Notice(to, "You are not a registered admin.")
		return
	}
	if c.pending {
		ib.Notice(to, "Wait for other operations to complete")
		return
	}
	c.pending = true

	ib.AddCallback("330", func(e *irc.Event) {
		c.auth = true
		ib.ClearCallback("330")
	})
	ib.AddCallback("318", func(e *irc.Event) {
		ib.ClearCallback("318")
		time.Sleep(1 * time.Second)
		if !c.auth {
			ib.Notice(to, "You must identify to nickserv in order to use this plugin.")
			c.pending = false
			return
		}
		c.processArgs(ib, to)
	})

	c.args = args

	ib.Whois(from)
}

// Help shows the help for the plugin.
func (c *Command) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Manages the bot configuration")
	ib.Privmsg(from, "`!config plugins` will give the state of all plugins and middlewares")
	ib.Privmsg(from, "`!config plugins (+|-)[p:|m:]pluginName will enable/disable a plugin")
	ib.Privmsg(from, "`!config admin` will list the admins")
	ib.Privmsg(from, "`!config admin [nicks...]` will add a space separated list of nicks to the admins list")
	ib.Privmsg(from, "`!config save` will save the current configuration to the config file used to load the initial config appended by \".new\"")
	ib.Privmsg(from, "`!config save truncate` will save the current configuration to the config file used to load the initial config")
}

// Start returns nil since it is a core plugin
func (c *Command) Start() error {
	c.toStop = new(modifier)
	c.toStart = new(modifier)
	c.auth = false
	c.pending = false
	c.args = nil
	return nil
}

// Stop returns nil since it is a core plugin
func (c *Command) Stop() error {
	return nil
}

// IsStarted returns always true since it is a core plugin
func (c *Command) IsStarted() bool {
	return true
}

func (c *Command) processArgs(ib *irc.Connection, to string) {
	cnf := configuration.Config
	if len(c.args) == 0 {
		c.Start()
		return
	}
	switch c.args[0] {
	case "save":
		if err := c.save(); err != nil {
			ib.Privmsg(to, "Error while saving configuration. Consult the logs.")
		} else {
			ib.Privmsg(to, "Configuration saved")
		}
		c.Start()
		return
	case "admins":
		if len(c.args) == 1 {
			ib.Privmsgf(to, "Admins : %v", cnf.Admins)
		} else {
			if c.admins() {
				ib.Privmsg(to, "Admins configuration changed")
			}
		}
		c.Start()
		return
	case "reset":
		// Can be done in a smarter way, we should copy the slices from configuration
		// and stop/start only those different from the default config
		plugins.Stop()
		configuration.Load()
		plugins.Start()
		ib.Privmsg(to, "Configuration was reseted")
		return
	case "plugins":
		if len(c.args) == 1 {
			list := list()
			ib.Privmsgf(to, "Commands     : %s", list[0])
			ib.Privmsgf(to, "Middlewares : %s", list[1])
			c.Start()
			return
		}
		if c.plugins() {
			ib.Privmsg(to, "Commands configuration changed")
		}
	}
	c.Start()
}

func (c *Command) save() error {
	if len(c.args) == 2 && c.args[1] == "truncate" {
		return configuration.Save(true)
	}
	return configuration.Save(false)
}

func (c *Command) admins() bool {
	var effective = false
	cnf := configuration.Config
	for _, i := range c.args[1:] {
		if strings.HasPrefix(i, "-") {
			cnf.Admins, _ = utils.RemoveStringInSlice(i[1:], cnf.Admins)
			effective = true
		} else if strings.HasPrefix(i, "+") && !utils.StringInSlice(i[1:], cnf.Admins) {
			cnf.Admins = append(cnf.Admins, i[1:])
			effective = true
		}
	}
	return effective
}

func (c *Command) plugins() bool {
	for _, i := range c.args[1:] {
		if strings.HasPrefix(i, "-") && len(i) > 1 {
			if i[1:3] == "m:" {
				c.toStop.Middlewares = append(c.toStop.Middlewares, i[3:])
			} else if i[1:3] == "p:" {
				c.toStop.Commands = append(c.toStop.Commands, i[3:])
			} else {
				c.toStop.Commands = append(c.toStop.Commands, i[1:])
				c.toStop.Middlewares = append(c.toStop.Middlewares, i[1:])
			}
		} else if strings.HasPrefix(i, "+") {
			if i[1:3] == "m:" {
				c.toStart.Middlewares = append(c.toStart.Middlewares, i[3:])
			} else if i[1:3] == "p:" {
				c.toStart.Commands = append(c.toStart.Commands, i[3:])
			} else {
				c.toStart.Commands = append(c.toStart.Commands, i[1:])
				c.toStart.Middlewares = append(c.toStart.Middlewares, i[1:])
			}
		}
	}
	return c.modify()
}

// modify executes the requested changes
func (c *Command) modify() bool {
	effective := false
	cnf := configuration.Config
	for _, n := range c.toStart.Commands {
		if _, ok := plugins.Commands[n]; ok {
			cnf.Commands = append(cnf.Commands, n)
			effective = true
		}
	}
	for _, n := range c.toStart.Middlewares {
		if _, ok := plugins.Middlewares[n]; ok {
			cnf.Middlewares = append(cnf.Middlewares, n)
			effective = true
		}
	}
	for _, n := range c.toStop.Commands {
		if _, ok := plugins.Commands[n]; ok && utils.StringInSlice(n, cnf.Commands) {
			cnf.Commands, _ = utils.RemoveStringInSlice(n, cnf.Commands)
			plugins.Commands[n].Stop()
			effective = true
		}
	}
	for _, n := range c.toStop.Middlewares {
		if _, ok := plugins.Middlewares[n]; ok && utils.StringInSlice(n, cnf.Middlewares) {
			cnf.Middlewares, _ = utils.RemoveStringInSlice(n, cnf.Middlewares)
			plugins.Middlewares[n].Stop()
			effective = true
		}
	}
	plugins.Start()
	return effective
}

// list returns a pair of strings listing all the plugins/middlewares and their state
func list() [2]string {
	var list [2]string
	for k, c := range plugins.Commands {
		if c.IsStarted() {
			list[0] += fmt.Sprintf(" \x0303+%s\x03", k)
		} else {
			list[0] += fmt.Sprintf(" \x0304-%s\x0F\x03", k)
		}
	}
	for k, c := range plugins.Middlewares {
		if c.IsStarted() {
			list[1] += fmt.Sprintf(" \x0303+%v\x03", k)
		} else {
			list[1] += fmt.Sprintf(" \x0304-%v\x0F\x03", k)
		}
	}
	return list
}

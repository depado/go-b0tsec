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
	pluginName = "config"
)

type modifier struct {
	Plugins     []string
	Middlewares []string
}

// Plugin is the help plugin
type Plugin struct {
	Pending bool
	Auth    bool
	ToStart *modifier
	ToStop  *modifier
}

// init initializes all the plugins and middlewares.
func init() {
	plugins.Plugins[pluginName] = new(Plugin)
}

// Get actually executes the command.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !utils.StringInSlice(from, configuration.Config.Admins) {
		ib.Privmsg(to, "You are not a registered admin.")
		return
	}
	if len(args) == 0 {
		list := List()
		ib.Privmsgf(to, "Plugins     : %s", list[0])
		ib.Privmsgf(to, "Middlewares : %s", list[1])
		return
	}

	if p.Pending {
		ib.Privmsg(to, "Wait for other operations to complete")
		return
	}
	p.Pending = true

	ib.AddCallback("330", func(e *irc.Event) {
		p.Auth = true
		ib.ClearCallback("330")
	})
	ib.AddCallback("318", func(e *irc.Event) {
		ib.ClearCallback("318")
		time.Sleep(1 * time.Second)
		if !p.Auth {
			p.ToStop = new(modifier)
			p.ToStart = new(modifier)
			p.Pending = false
			ib.Privmsg(to, "You must identify to nickserv in order to use this plugin.")
			return
		}
		p.Modify()
	})

	p.processArgs(args)

	ib.Whois(from)
}

// Help shows the help for the plugin.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Manages the bot configuration")
}

// Start returns nil since it is a core plugin
func (p *Plugin) Start() error {
	p.ToStop = new(modifier)
	p.ToStart = new(modifier)
	return nil
}

// Stop returns nil since it is a core plugin
func (p *Plugin) Stop() error {
	return nil
}

// IsStarted returns always true since it is a core plugin
func (p *Plugin) IsStarted() bool {
	return true
}

func (p *Plugin) processArgs(args []string) {
	for _, i := range args {
		if strings.HasPrefix(i, "-") {
			if i[1:2] == "m:" {
				p.ToStop.Middlewares = append(p.ToStop.Middlewares, i[3:])
			} else if i[1:2] == "p:" {
				p.ToStop.Plugins = append(p.ToStop.Plugins, i[3:])
			} else {
				p.ToStop.Plugins = append(p.ToStop.Plugins, i[1:])
				p.ToStop.Middlewares = append(p.ToStop.Middlewares, i[1:])
			}
		} else if strings.HasPrefix(i, "+") {
			if i[1:2] == "m:" {
				p.ToStart.Middlewares = append(p.ToStart.Middlewares, i[3:])
			} else if i[1:2] == "p:" {
				p.ToStart.Plugins = append(p.ToStart.Plugins, i[3:])
			} else {
				p.ToStart.Plugins = append(p.ToStart.Plugins, i[1:])
				p.ToStart.Middlewares = append(p.ToStart.Middlewares, i[1:])
			}
		}
	}
}

// Modify executes the requested changes
func (p *Plugin) Modify() {
	cnf := configuration.Config
	for _, n := range p.ToStart.Plugins {
		if !utils.StringInSlice(n, cnf.Plugins) {
			if _, ok := plugins.Plugins[n]; ok {
				cnf.Plugins = append(cnf.Plugins, n)
			}
		}
	}
	for _, n := range p.ToStart.Middlewares {
		if !utils.StringInSlice(n, cnf.Middlewares) {
			if _, ok := plugins.Middlewares[n]; ok {
				cnf.Middlewares = append(cnf.Middlewares, n)
			}
		}
	}
	var removed bool
	for _, n := range p.ToStop.Plugins {
		cnf.Plugins, removed = utils.RemoveStringInSlice(n, cnf.Plugins)
		if removed {
			plugins.Plugins[n].Stop()
		}
	}

	for _, n := range p.ToStop.Middlewares {
		cnf.Middlewares, removed = utils.RemoveStringInSlice(n, cnf.Middlewares)
		if removed {
			plugins.Middlewares[n].Stop()
		}
	}
	plugins.Start()
	p.ToStart = new(modifier)
	p.ToStop = new(modifier)
	p.Auth = false
	p.Pending = false
}

// List returns a pair of strings listing all the plugins/middlewares and their state
func List() [2]string {
	var list [2]string
	for k, p := range plugins.Plugins {
		if p.IsStarted() {
			list[0] += fmt.Sprintf(" \x0303+%s\x03", k)
		} else {
			list[0] += fmt.Sprintf(" \x0304-%s\x0F\x03", k)
		}
	}
	for k, p := range plugins.Middlewares {
		if p.IsStarted() {
			list[1] += fmt.Sprintf(" \x0303+%s\x03", k)
		} else {
			list[1] += fmt.Sprintf(" \x0304-%s\x0F\x03", k)
		}
	}
	return list
}

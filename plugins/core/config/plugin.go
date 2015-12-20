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
	pending bool
	auth    bool
	toStart *modifier
	toStop  *modifier
	args    []string
}

// init initializes all the plugins and middlewares.
func init() {
	plugins.Plugins[pluginName] = new(Plugin)
}

// Get actually executes the command.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if to == configuration.Config.BotName {
		to = from
	}
	if !utils.StringInSlice(from, configuration.Config.Admins) {
		ib.Notice(to, "You are not a registered admin.")
		return
	}

	if p.pending {
		ib.Notice(to, "Wait for other operations to complete")
		return
	}
	p.pending = true

	ib.AddCallback("330", func(e *irc.Event) {
		p.auth = true
		ib.ClearCallback("330")
	})
	ib.AddCallback("318", func(e *irc.Event) {
		ib.ClearCallback("318")
		time.Sleep(1 * time.Second)
		if !p.auth {
			ib.Notice(to, "You must identify to nickserv in order to use this plugin.")
			p.pending = false
			return
		}
		p.processArgs(ib, to)
	})

	p.args = args

	ib.Whois(from)
}

// Help shows the help for the plugin.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Manages the bot configuration")
	ib.Privmsg(from, "`!config plugins` will give the state of all plugins and middlewares")
	ib.Privmsg(from, "`!config plugins (+|-)[p:|m:]pluginName wil enable/disable a plugin")
	ib.Privmsg(from, "`!config admin` will list the admins")
	ib.Privmsg(from, "`!config admin [nicks...]` will add a space separated list of nicks to the admins list")
	ib.Privmsg(from, "`!config save` will save the current configuration to the config file used to load the initial config appended by \".new\"")
	ib.Privmsg(from, "`!config save truncate` will save the current configuration to the config file used to load the initial config")
}

// Start returns nil since it is a core plugin
func (p *Plugin) Start() error {
	p.toStop = new(modifier)
	p.toStart = new(modifier)
	p.auth = false
	p.pending = false
	p.args = nil
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

func (p *Plugin) processArgs(ib *irc.Connection, to string) {
	cnf := configuration.Config
	switch p.args[0] {
	case "save":
		p.save()
		p.Start()
		return
	case "admins":
		if len(p.args) == 1 {
			ib.Privmsgf(to, "Admins : %v", cnf.Admins)
		} else {
			p.admins()
		}
		p.Start()
		return
	case "reset":
		// Can be done in a smarter way, we should copy the slices from configuration
		// and stop/start only those different from the default config
		plugins.Stop()
		configuration.Load()
		plugins.Start()
		return
	case "plugins":
		if len(p.args) == 1 {
			list := list()
			ib.Privmsgf(to, "Plugins     : %s", list[0])
			ib.Privmsgf(to, "Middlewares : %s", list[1])
			p.Start()
			return
		}
		p.plugins()
	}
}

func (p *Plugin) save() {
	if len(p.args) == 2 && p.args[1] == "truncate" {
		configuration.Save(true)
	} else {
		configuration.Save(false)
	}
}

func (p *Plugin) admins() {
	cnf := configuration.Config
	for _, i := range p.args[1:] {
		if strings.HasPrefix(i, "-") {
			cnf.Admins, _ = utils.RemoveStringInSlice(i[1:], cnf.Admins)
		} else if !utils.StringInSlice(i, cnf.Admins) {
			cnf.Admins = append(cnf.Admins, i)
		}
	}
}

func (p *Plugin) plugins() {

	for _, i := range p.args[1:] {
		if strings.HasPrefix(i, "-") && len(i) > 1 {
			if i[1:3] == "m:" {
				p.toStop.Middlewares = append(p.toStop.Middlewares, i[3:])
			} else if i[1:3] == "p:" {
				p.toStop.Plugins = append(p.toStop.Plugins, i[3:])
			} else {
				p.toStop.Plugins = append(p.toStop.Plugins, i[1:])
				p.toStop.Middlewares = append(p.toStop.Middlewares, i[1:])
			}
		} else if strings.HasPrefix(i, "+") {
			if i[1:3] == "m:" {
				p.toStart.Middlewares = append(p.toStart.Middlewares, i[3:])
			} else if i[1:3] == "p:" {
				p.toStart.Plugins = append(p.toStart.Plugins, i[3:])
			} else {
				p.toStart.Plugins = append(p.toStart.Plugins, i[1:])
				p.toStart.Middlewares = append(p.toStart.Middlewares, i[1:])
			}
		}
	}
	p.modify()
}

// modify executes the requested changes
func (p *Plugin) modify() bool {
	effective := false
	cnf := configuration.Config
	for _, n := range p.toStart.Plugins {
		if _, ok := plugins.Plugins[n]; ok {
			cnf.Plugins = append(cnf.Plugins, n)
			effective = true
		}
	}
	for _, n := range p.toStart.Middlewares {
		if _, ok := plugins.Middlewares[n]; ok {
			cnf.Middlewares = append(cnf.Middlewares, n)
			effective = true
		}
	}
	for _, n := range p.toStop.Plugins {
		if _, ok := plugins.Plugins[n]; ok && utils.StringInSlice(n, cnf.Plugins) {
			cnf.Plugins, _ = utils.RemoveStringInSlice(n, cnf.Plugins)
			plugins.Plugins[n].Stop()
			effective = true
		}
	}
	for _, n := range p.toStop.Middlewares {
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
	for k, p := range plugins.Plugins {
		if p.IsStarted() {
			list[0] += fmt.Sprintf(" \x0303+%s\x03", k)
		} else {
			list[0] += fmt.Sprintf(" \x0304-%s\x0F\x03", k)
		}
	}
	for k, p := range plugins.Middlewares {
		if p.IsStarted() {
			list[1] += fmt.Sprintf(" \x0303+%v\x03", k)
		} else {
			list[1] += fmt.Sprintf(" \x0304-%v\x0F\x03", k)
		}
	}
	return list
}

package afk

import (
	"strings"
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "afk"
)

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct {
	Started bool
}

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "Tell the world you're afk, for a reason. Or not.")
	ib.Privmsg(from, "Example : !afk reason.")
}

// Get is the actual call to your plugin.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
		return
	}
	reason := ""
	if len(args) > 0 {
		reason = strings.Join(args, " ")
	}
	Map[from] = Data{time.Now(), reason}
	if reason != "" {
		ib.Privmsgf(configuration.Config.Channel, "%v is afk : %v", from, reason)
	} else {
		ib.Privmsgf(configuration.Config.Channel, "%v is afk.", from)
	}
}

// Start starts the plugin and returns any occured error, nil otherwise
func (p *Plugin) Start() error {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		p.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occured error, nil otherwise
func (p *Plugin) Stop() error {
	p.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (p *Plugin) IsStarted() bool {
	return p.Started
}

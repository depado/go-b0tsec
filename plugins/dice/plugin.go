package dice

import (
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "dice"
)

var re = regexp.MustCompile(`(\d+)d(\d+)`)

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct{}

func init() {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		plugins.Plugins[pluginCommand] = new(Plugin)
	}
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Throws x n-faced dice(s) in the form of 'xdn'")
	ib.Privmsg(from, "Example with 2 20-faced dices : !dice 2d20")
	ib.Privmsg(from, "Limit : 25d1000")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		fmtstr := "%v"
		if len(args) > 1 {
			fmtstr = "%v - " + strings.Join(args[1:], " ")
		}
		rs := re.FindAllStringSubmatch(args[0], -1)
		if len(rs) > 0 {
			t, err := strconv.Atoi(rs[0][1])
			if err != nil {
				return
			}
			dt, err := strconv.Atoi(rs[0][2])
			if err != nil {
				return
			}
			if dt == 0 || dt > 1000 || t > 25 {
				return
			}
			if to == configuration.Config.BotName {
				ib.Privmsgf(from, fmtstr, throw(t, dt))
			} else {
				ib.Privmsgf(to, "%v: "+fmtstr, from, throw(t, dt))
			}
		}
	}
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}

func throw(times int, dice int) string {
	tot := 0
	ds := make([]string, times)
	for i := 0; i < times; i++ {
		c := rand.Intn(dice)
		tot += c + 1
		ds[i] = strconv.Itoa(c + 1)
	}
	if times > 1 {
		return strings.Join(ds, " + ") + " = " + strconv.Itoa(tot)
	}
	return strconv.Itoa(tot)
}

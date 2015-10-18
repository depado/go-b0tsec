package karma

import (
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"
)

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct {
	Karma  map[string]int
	Action map[string]time.Time
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    Karma stuff")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if len(args) > 0 {
		switch args[0] {
		case "+", "-":
			if len(args) > 1 {
				if from != args[1] {
					if val, ok := p.Action[from]; ok {
						if time.Since(val) < 30*time.Minute {
							ib.Notice(from, "Please wait 30 minutes between each karma operation.")
							return
						}
					}
					p.Action[from] = time.Now()
					c := 0
					if val, ok := p.Karma[args[1]]; ok {
						c = val
					}
					if args[0] == "+" {
						p.Karma[args[1]] = c + 1
						ib.Privmsgf(configuration.Config.Channel, "Someone gave a karma point to %v, total %v", args[1], c+1)
					} else {
						p.Karma[args[1]] = c - 1
						ib.Privmsgf(configuration.Config.Channel, "Someone took a karma point from %v, total %v", args[1], c-1)
					}
				} else {
					ib.Notice(from, "Can't add or remove points to yourself.")
					return
				}
			} else {
				ib.Notice(from, "You need to give a nickname to operate on.")
			}
		case "=":
			if len(args) > 1 {
				for _, n := range args[1:] {
					if val, ok := p.Karma[n]; ok {
						ib.Privmsgf(to, "%v has %v point(s).", n, val)
					} else {
						ib.Privmsgf(to, "I don't have records on %v.", n)
					}
				}
			} else {
				ib.Notice(from, "Need at least a nickname.")
			}
		}
	}
}

// New initializes new plugin
func New() Plugin {
	return Plugin{make(map[string]int), make(map[string]time.Time)}
}

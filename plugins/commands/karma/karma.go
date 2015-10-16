package example

import "github.com/thoj/go-ircevent"

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct{}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    This command does this and that. For example this. And that.")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	ib.Privmsg(to, "This command doesn't do much, but I'm glad you tried it out.")
}

package example

import "github.com/thoj/go-ircevent"

/*
   This is a plugin example. This is not intended to be activated.
   If you want to activate your plugin in the bot, go to the file named plugins.go and execute the following call :
   RegisterCommand("command", new(yourpackage.Plugin))
   This will insert your command into the plugin map so that it can be called like that : !command
   Once your plugin is activated, it will also display the help when calling !help or !help command
*/

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
// It also makes sure you don't have two plugins coming from the same package.
// This struct can also be used to save a state that will persist across multiple commands.
// For a more complex example, see the urban plugin.
type Plugin struct{}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
// As this command is intented to display long text, it will send only to the issuer of the command.
// You can however display the help anywhere else if you wish to (using configuration.Config.Channel for example).
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "    This command does this and that. For example this. And that.")
}

// Get is the actual call to your plugin.
// from: The nickname which executed the command. Useful to highlight said person or send a private message.
// to: The channel in which the command has been executed. You can determine if the command has been sent directly to the bot in a private message for example.
// args: An array of strings, representing everything after the command call.
// Note : If you intend to modify the state of the plugin, (p Plugin) should then become (p *Plugin)
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	ib.Privmsg(to, "This command doesn't do much, but I'm glad you tried it out.")
}

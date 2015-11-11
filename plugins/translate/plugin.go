package translate

import (
	"log"
	"net/http"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/translate/v2"
)

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct{}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "This command will never work due to Google being huge assholes.")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if to == configuration.Config.BotName {
		to = from
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: configuration.Config.GoogleAPIKey},
	}
	service, err := translate.New(client)
	if err != nil {
		log.Printf("Error creating translate service : %s\n", err)
		return
	}
	if len(args) > 2 && args[len(args)-2] == ">" {
		lang := args[len(args)-1]
		trservice := translate.NewTranslationsService(service)
		resp, err := trservice.List([]string{strings.Join(args[:len(args)-2], " ")}, lang).Do()
		if err != nil {
			ib.Privmsg(to, "Are you trying to use that ? You can't. Give me 5$ and I'll activate this for you.")
			return
		}
		for _, tr := range resp.Translations {
			log.Println(tr)
		}
	}
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}

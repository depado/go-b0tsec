package translate

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

// Yandex struct represents the response of the Yandex translate API.
type Yandex struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

var translateEndpoint string

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct{}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Translate things from one language to another.")
	ib.Privmsg(from, "Example (english to french) : !translate flying saucer > fr")
	ib.Privmsg(from, "Example (french to english) : !translate soucoupe volante > en")
}

// Get is the actual call to your plugin.
func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if to == configuration.Config.BotName {
		to = from
	}
	if len(args) > 2 && args[len(args)-2] == ">" {
		lang := args[len(args)-1]
		q := url.QueryEscape(strings.Join(args[:len(args)-2], " "))
		endpoint := fmt.Sprintf(translateEndpoint, lang, q)
		yr := Yandex{}
		if err := utils.FetchURL(endpoint, &yr); err != nil {
			log.Println(err)
			return
		}
		ib.Privmsgf(to, "%v: \x0314[%v]\x0F\x03 %v", from, yr.Lang, yr.Text[0])
	}
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	translateEndpoint = "https://translate.yandex.net/api/v1.5/tr.json/translate?key=" + configuration.Config.YandexTrnslKey + "&lang=%s&text=%s"
	return new(Plugin)
}

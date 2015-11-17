package define

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

var dictionnaryEndpoint string

// YandexDict struct holds the response of a call to the Yandex dictionnary API.
type YandexDict struct {
	Head struct {
	} `json:"head"`
	Def []struct {
		Text string `json:"text"`
		Pos  string `json:"pos"`
		Tr   []struct {
			Text string `json:"text"`
			Pos  string `json:"pos"`
			Syn  []struct {
				Text string `json:"text"`
			} `json:"syn"`
			Mean []struct {
				Text string `json:"text"`
			} `json:"mean"`
			Ex []struct {
				Text string `json:"text"`
				Tr   []struct {
					Text string `json:"text"`
				} `json:"tr"`
			} `json:"ex"`
		} `json:"tr"`
	} `json:"def"`
}

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
	if len(args) > 2 && args[len(args)-2] == ">" {
		lang := fmt.Sprintf("%v-en", args[len(args)-1])
		q := url.QueryEscape(strings.Join(args[:len(args)-2], " "))
		endpoint := fmt.Sprintf(dictionnaryEndpoint, lang, q)
		yr := YandexDict{}
		if err := utils.FetchURL(endpoint, &yr); err != nil {
			log.Println(err)
			return
		}
		for _, def := range yr.Def {
			mean := ""
			for _, tr := range def.Tr {
				log.Println(tr)
				for _, m := range tr.Mean {
					if m.Text != "" {
						mean = ": " + m.Text
					}
				}
			}
			ib.Privmsgf(to, "%v - %v %v", def.Text, def.Pos, mean)
		}
	}
}

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	dictionnaryEndpoint = "https://dictionary.yandex.net/api/v1/dicservice.json/lookup?key=" + configuration.Config.YandexDictKey + "&lang=%s&text=%s"
	return new(Plugin)
}

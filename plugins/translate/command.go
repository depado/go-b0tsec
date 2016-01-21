package translate

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	command = "translate"
)

var translateEndpoint = "https://translate.yandex.net/api/v1.5/tr.json/translate?key=%s&lang=%s&text=%s"

// Yandex struct represents the response of the Yandex translate API.
type Yandex struct {
	Code int      `json:"code"`
	Lang string   `json:"lang"`
	Text []string `json:"text"`
}

// Command is the plugin struct. It will be exposed as packagename.Command to keep the API stable and friendly.
type Command struct {
	Started bool
}

func init() {
	plugins.Commands[command] = new(Command)
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Translate things from one language to another.")
	ib.Privmsg(from, "Example (english to french) : !translate flying saucer > fr")
	ib.Privmsg(from, "Example (french to english) : !translate soucoupe volante > en")
}

// Get is the actual call to your plugin.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	cnf := configuration.Config
	if !c.Started {
		return
	}
	if to == cnf.BotName {
		to = from
	}
	if len(args) > 2 && args[len(args)-2] == ">" {
		lang := args[len(args)-1]
		q := url.QueryEscape(strings.Join(args[:len(args)-2], " "))
		endpoint := fmt.Sprintf(translateEndpoint, cnf.YandexTrnslKey, lang, q)
		yr := Yandex{}
		if err := utils.FetchURL(endpoint, &yr); err != nil {
			log.Println(err)
			return
		}
		if len(yr.Text) > 0 {
			ib.Privmsgf(to, "%v: \x0314[%v]\x0F\x03 %v", from, yr.Lang, yr.Text[0])
		} else {
			ib.Privmsgf(to, "%v: \x0314Unrecognised language.\x0F\x03", from)
		}
	}
}

// Start starts the plugin and returns any occurred error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(command, configuration.Config.Commands) {
		c.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occurred error, nil otherwise
func (c *Command) Stop() error {
	c.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (c *Command) IsStarted() bool {
	return c.Started
}

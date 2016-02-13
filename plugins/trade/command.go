package trade

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thoj/go-ircevent"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
)

const (
	command = "trade"
	api     = "https://www.cryptonator.com/api/full/%s-%s"
)

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
	ib.Privmsg(from, "Displays current crypto-currency trade market")
	ib.Privmsg(from, "Usage : !trade [ammount] <from> <to> [--market=all|market_name] [--nomarket] [--sort=price|volume]")
}

// Get is the actual call to your plugin.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	var err error
	var opt options
	if to == configuration.Config.BotName {
		to = from
	}
	if opt, err = parseOptions(args); err != nil {
		ib.Privmsg(to, err.Error())
		return
	}
	trade := Cryptonator{}
	if err = utils.FetchURL(fmt.Sprintf(api, opt.Source, opt.Dest), &trade); err == nil {
		rate, err := strconv.ParseFloat(trade.Ticker.Price, 64)
		if err != nil {
			ib.Privmsgf(to, "Couldn't get the price of %s or %s", opt.Source, opt.Dest)
			return
		}
		if opt.Ammount == 1 {
			ib.Privmsgf(to, "%v %s = %f %s", opt.Ammount, trade.Ticker.Base, float64(opt.Ammount)*rate, trade.Ticker.Target)
		} else {
			ib.Privmsgf(to, "%v %s = %f %s (1 %s = %f %s)", opt.Ammount, trade.Ticker.Base, float64(opt.Ammount)*rate, trade.Ticker.Target, trade.Ticker.Base, rate, trade.Ticker.Target)
		}
		if opt.Market != "" {
			if opt.Market == "all" {
				for _, m := range trade.Ticker.Markets {
					ib.Privmsgf(to, "[%s] %s %s [%f]", m.Market, m.Price, trade.Ticker.Target, m.Volume)
				}
			} else {
				for _, m := range trade.Ticker.Markets {
					if strings.ToLower(m.Market) == strings.ToLower(opt.Market) {
						ib.Privmsgf(to, "[%s] %s %s [%f]", m.Market, m.Price, trade.Ticker.Target, m.Volume)
					}
				}
			}
		}
	} else {
		fmt.Println(err)
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

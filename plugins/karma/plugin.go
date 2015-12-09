package karma

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/database"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	pluginCommand = "karma"
	bucketName    = "karma"
	mainKey       = "main"
)

// Data is the struct that contains the data about the karma intented to be stored somewhere.
type Data struct {
	Karma map[string]int
}

// Pair is used to sort the map
type Pair struct {
	Key   string
	Value int
}

// PairList is a list of Pair
type PairList []Pair

// Plugin is the plugin struct. It will be exposed as packagename.Plugin to keep the API stable and friendly.
type Plugin struct {
	Started bool
	Data
	Action map[string]time.Time
}

func init() {
	plugins.Plugins[pluginCommand] = new(Plugin)
}

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Encode encodes a chain to json.
func (d Data) Encode() ([]byte, error) {
	enc, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

// Decode decodes json to Chain
func (d *Data) Decode(data []byte) error {
	if err := json.Unmarshal(data, d); err != nil {
		return err
	}
	return nil
}

// Save saves the Data
func (d Data) Save() error {
	return database.BotStorage.Save(bucketName, mainKey, &d)
}

// CanModify checks if the args are correct and if the user can modify a karma
func (p *Plugin) CanModify(from string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("You need to give a nickname to operate on.")
	}
	if from == args[1] {
		return fmt.Errorf("Can't add or remove points to yourself.")
	}
	if val, ok := p.Action[from]; ok {
		if time.Since(val) < 1*time.Minute {
			return fmt.Errorf("Please wait 1 minute between each karma operation.")
		}
	}
	return nil
}

// ModifyKarma modifies the karma value of a user and send msg to notify the new count
func (p *Plugin) ModifyKarma(ib *irc.Connection, from string, to string, args []string) {
	if err := p.CanModify(from, args); err != nil {
		ib.Notice(from, err.Error())
		return
	}
	p.Action[from] = time.Now()
	c := 0
	if val, ok := p.Karma[args[1]]; ok {
		c = val
	}
	if args[0] == ">" {
		p.Karma[args[1]] = c + 1
		ib.Privmsgf(configuration.Config.Channel, "Someone gave a karma point to %v, total %v", args[1], c+1)
	} else {
		p.Karma[args[1]] = c - 1
		ib.Privmsgf(configuration.Config.Channel, "Someone took a karma point from %v, total %v", args[1], c-1)
	}
	if len(args) > 2 {
		ib.Privmsgf(configuration.Config.Channel, "Reason : %s", strings.Join(args[2:], " "))
	}
	if err := p.Data.Save(); err != nil {
		log.Println(err)
	}
}

// GetKarma gets karma value of all
func (p *Plugin) GetKarma(ib *irc.Connection, from string, to string, nicks []string) {
	if len(nicks) < 1 {
		ib.Notice(from, "You need to give a nickname to operate on.")
		return
	}
	nrec := false
	pl := make(PairList, 0)
	for _, n := range nicks {
		if val, ok := p.Karma[n]; ok {
			pl = append(pl, Pair{n, val})
		} else {
			nrec = true
		}
	}
	sort.Sort(sort.Reverse(pl))
	for _, v := range pl {
		ib.Privmsgf(to, "%v has %v point(s)", v.Key, v.Value)
	}
	if nrec {
		ib.Privmsg(to, "No record for the others.")
	}
}

// Help must send some help about what the command actually does and how to call it if there are any optional arguments.
func (p *Plugin) Help(ib *irc.Connection, from string) {
	if !p.Started {
		return
	}
	ib.Privmsg(from, "Allows to add/remove/see karma points to/from/of a person.")
	ib.Privmsg(from, "Add    : !karma > nickname [optional reason]")
	ib.Privmsg(from, "Remove : !karma < nickname [optional reason]")
	ib.Privmsg(from, "See    : !karma = nickname1 [nickname2, nickname3, ...]")
}

// Get is the actual call to your plugin.
func (p *Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
	if !p.Started {
		return
	}
	if len(args) > 0 {
		switch args[0] {
		case "<", ">":
			p.ModifyKarma(ib, from, to, args)
			return
		case "=":
			p.GetKarma(ib, from, to, args[1:])
			return
		}
	}
}

// Start starts the plugin and returns any occured error, nil otherwise
func (p *Plugin) Start() error {
	if utils.StringInSlice(pluginCommand, configuration.Config.Plugins) {
		if err := database.BotStorage.CreateBucket(bucketName); err != nil {
			log.Fatalf("Error while creating bucket for the Karma plugin : %s", err)
		}
		d := Data{make(map[string]int)}
		database.BotStorage.Get(bucketName, mainKey, &d)

		plugins.Plugins[pluginCommand] = &Plugin{true, d, make(map[string]time.Time)}
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

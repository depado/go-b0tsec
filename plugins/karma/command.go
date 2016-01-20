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
	command    = "karma"
	bucketName = "karma"
	mainKey    = "main"
)

// Data is the struct that contains the data about the karma intended to be stored somewhere.
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

// Command is the plugin struct. It will be exposed as packagename.Command to keep the API stable and friendly.
type Command struct {
	Started bool
	Data
	Action map[string]time.Time
}

func init() {
	plugins.Commands[command] = new(Command)
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
func (c *Command) CanModify(from string, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("You need to give a nickname to operate on.")
	}
	if from == args[1] {
		return fmt.Errorf("Can't add or remove points to yourself.")
	}
	if val, ok := c.Action[from]; ok {
		if time.Since(val) < 1*time.Minute {
			return fmt.Errorf("Please wait 1 minute between each karma operation.")
		}
	}
	return nil
}

// ModifyKarma modifies the karma value of a user and send msg to notify the new count
func (c *Command) ModifyKarma(ib *irc.Connection, from string, to string, args []string) {
	if err := c.CanModify(from, args); err != nil {
		ib.Notice(from, err.Error())
		return
	}
	c.Action[from] = time.Now()
	current := 0
	if val, ok := c.Karma[args[1]]; ok {
		current = val
	}
	if args[0] == ">" {
		c.Karma[args[1]] = current + 1
		ib.Privmsgf(configuration.Config.Channel, "Someone gave a karma point to %v, total %v", args[1], current+1)
	} else {
		c.Karma[args[1]] = current - 1
		ib.Privmsgf(configuration.Config.Channel, "Someone took a karma point from %v, total %v", args[1], current-1)
	}
	if len(args) > 2 {
		ib.Privmsgf(configuration.Config.Channel, "Reason : %s", strings.Join(args[2:], " "))
	}
	if err := c.Data.Save(); err != nil {
		log.Println(err)
	}
}

// GetKarma gets karma value of all
func (c *Command) GetKarma(ib *irc.Connection, from string, to string, nicks []string) {
	if len(nicks) < 1 {
		ib.Notice(from, "You need to give a nickname to operate on.")
		return
	}
	nrec := false
	pl := make(PairList, 0)
	for _, n := range nicks {
		if val, ok := c.Karma[n]; ok {
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
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Allows to add/remove/see karma points to/from/of a person.")
	ib.Privmsg(from, "Add    : !karma > nickname [optional reason]")
	ib.Privmsg(from, "Remove : !karma < nickname [optional reason]")
	ib.Privmsg(from, "See    : !karma = nickname1 [nickname2, nickname3, ...]")
}

// Get is the actual call to your plugin.
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		switch args[0] {
		case "<", ">":
			c.ModifyKarma(ib, from, to, args)
			return
		case "=":
			c.GetKarma(ib, from, to, args[1:])
			return
		}
	}
}

// Start starts the plugin and returns any occurred error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(command, configuration.Config.Commands) {
		if err := database.BotStorage.CreateBucket(bucketName); err != nil {
			log.Fatalf("Error while creating bucket for the Karma plugin : %s", err)
		}
		d := Data{make(map[string]int)}
		database.BotStorage.Get(bucketName, mainKey, &d)

		plugins.Commands[command] = &Command{true, d, make(map[string]time.Time)}
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

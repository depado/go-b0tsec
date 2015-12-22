package plugins

import (
	"log"

	"github.com/thoj/go-ircevent"
)

// Command represents a single plugin. The Get method is use to send things.
type Command interface {
	Get(*irc.Connection, string, string, []string)
	Help(*irc.Connection, string)
	Start() error
	Stop() error
	IsStarted() bool
}

// Middleware represents a single plugin. The Get method is use to send things.
type Middleware interface {
	Get(*irc.Connection, string, string, string)
	Start() error
	Stop() error
	IsStarted() bool
}

// Plugins is the map structure of all configured plugins
var Commands = map[string]Command{}

// Middlewares is the slice of all configured middlewares Get() func
var Middlewares = map[string]Middleware{}

// ListCommands returns a list of the started plugins
func ListCommands() []string {
	var list []string
	for k, p := range Commands {
		if p.IsStarted() {
			list = append(list, k)
		}
	}
	return list
}

// Stop calls the Stop method of each registered plugins
func Stop() {
	for _, k := range Middlewares {
		if err := k.Stop(); err != nil {
			log.Printf("Error closing middlewares : %v", err)
		}
	}

	for _, k := range Commands {
		if err := k.Stop(); err != nil {
			log.Printf("Error closing plugins : %v", err)
		}
	}
}

// Start calls the Start method of each registered plugin
func Start() {
	for _, k := range Middlewares {
		if err := k.Start(); err != nil {
			log.Printf("Error starting middlewares : %v", err)
		}
	}
	for p, k := range Commands {
		if err := k.Start(); err != nil {
			log.Printf("Error starting plugin %s : %s\n", p, err)
		}
	}
}

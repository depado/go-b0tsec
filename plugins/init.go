package plugins

import (
	"log"

	"github.com/thoj/go-ircevent"
)

// Plugin represents a single plugin. The Get method is use to send things.
type Plugin interface {
	//IsStarted() bool
	Get(*irc.Connection, string, string, []string)
	Help(*irc.Connection, string)
	//Stop() error
}

// Middleware represents a single plugin. The Get method is use to send things.
type Middleware interface {
	Get(*irc.Connection, string, string, string)
	Start() error
	Stop() error
	IsStarted() bool
}

// Plugins is the map structure of all configured plugins
var Plugins = map[string]Plugin{}

// Middlewares is the slice of all configured middlewares Get() func
var Middlewares = []Middleware{}

// Stop calls the Stop method of each registered middleware
func Stop() {
	for _, k := range Middlewares {
		if err := k.Stop(); err != nil {
			log.Println("Error closing middlewares : %s", err)
		}
	}
}

// Start calls the Start method of each registered middleware
func Start() {
	for _, k := range Middlewares {
		if err := k.Start(); err != nil {
			log.Println("Error starting middlewares : %s", err)
		}
	}
}

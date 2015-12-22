package youtube

import (
	"log"
	"net/http"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"

	"google.golang.org/api/googleapi/transport"
	yt "google.golang.org/api/youtube/v3"
)

const (
	pluginName    = "youtube"
	pluginCommand = "yt"
)

// Command is the plugin struct. It will be exposed as packagename.Command to keep the API stable and friendly.
type Command struct {
	Started bool
}

func init() {
	plugins.Commands[pluginCommand] = new(Command)
}

// Help shows a help message for this command.
func (c *Command) Help(ib *irc.Connection, from string) {
	if !c.Started {
		return
	}
	ib.Privmsg(from, "Search directly on YouTube.")
	ib.Privmsg(from, "Example : !yt Funny cat")
}

// Get actually sends the data
func (c *Command) Get(ib *irc.Connection, from string, to string, args []string) {
	if !c.Started {
		return
	}
	if len(args) > 0 {
		client := &http.Client{
			Transport: &transport.APIKey{Key: configuration.Config.GoogleAPIKey},
		}

		service, err := yt.New(client)
		if err != nil {
			log.Printf("Error creating new YouTube client: %v", err)
		}

		// Make the API call to YouTube.
		call := service.Search.List("id").Q(strings.Join(args, " ")).
			SafeSearch("none").Type("video").MaxResults(1)
		response, err := call.Do()
		if err != nil {
			log.Printf("Error making search API call: %v", err)
		}

		for _, i := range response.Items {
			videos, err := service.Videos.List("snippet, statistics, contentDetails").Id(i.Id.VideoId).Do()
			if err != nil {
				log.Println(err)
				return
			}
			for _, val := range videos.Items {
				ib.Privmsgf(to, "\u0002%v\u000F : %s",
					"https://youtu.be/"+i.Id.VideoId,
					FormatOutput(val))
			}
		}
	}
}

// Start starts the plugin and returns any occured error, nil otherwise
func (c *Command) Start() error {
	if utils.StringInSlice(pluginCommand, configuration.Config.Commands) {
		c.Started = true
	}
	return nil
}

// Stop stops the plugin and returns any occured error, nil otherwise
func (c *Command) Stop() error {
	c.Started = false
	return nil
}

// IsStarted returns the state of the plugin
func (c *Command) IsStarted() bool {
	return c.Started
}

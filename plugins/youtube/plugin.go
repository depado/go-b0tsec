package youtube

import (
	"log"
	"net/http"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"

	"google.golang.org/api/googleapi/transport"
	yt "google.golang.org/api/youtube/v3"
)

func (p Plugin) Help(ib *irc.Connection, from string) {
	ib.Privmsg(from, "Search directly on YouTube.")
	ib.Privmsg(from, "Example : !yt Funny cat")
}

func (p Plugin) Get(ib *irc.Connection, from string, to string, args []string) {
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

// NewPlugin returns a new plugin
func NewPlugin() *Plugin {
	return new(Plugin)
}

package youtube

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"

	"google.golang.org/api/googleapi/transport"
	yt "google.golang.org/api/youtube/v3"
)

var ytre = regexp.MustCompile(`(?:https?://)?(?:(?:www\.)?youtube\.com/watch\?.*v=|youtu\.be/)([^&?]{11})`)

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: configuration.Config.GoogleAPIKey},
	}
	for _, bit := range strings.Fields(message) {
		rs := ytre.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			service, err := yt.New(client)
			if err != nil {
				log.Printf("Error creating new YouTube client: %v\n", err)
				return
			}
			response, err := service.Videos.List("snippet, statistics, contentDetails").Id(rs[0][1]).Do()
			if err != nil {
				log.Println(err)
				return
			}
			for _, val := range response.Items {
				ib.Privmsgf(to, FormatOutput(val))
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

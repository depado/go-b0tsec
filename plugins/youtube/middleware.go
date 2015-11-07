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

var ytre = regexp.MustCompile(`(?:https?://)?(?:www.)?youtube.com/watch\?.*v=([^&]*)`)

// Middleware is the github middleware
type Middleware struct{}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: configuration.Config.YoutubeKey},
	}
	for _, bit := range strings.Fields(message) {
		rs := ytre.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			service, err := yt.New(client)
			if err != nil {
				log.Fatalf("Error creating new YouTube client: %v", err)
			}
			response, err := service.Videos.List("snippet, statistics, contentDetails").Id(rs[0][1]).Do()
			if err != nil {
				log.Println(err)
			}
			for _, val := range response.Items {
				t := strings.Replace(val.ContentDetails.Duration[2:len(val.ContentDetails.Duration)-1], "M", ":", -1)
				t = strings.Replace(t, "H", ":", -1)
				if err != nil {
					log.Println(err)
				}
				ib.Privmsgf(to, "%v [\x0303%v\x03 | \x0304%v\x0F\x03] (%v) (%v)",
					val.Snippet.Title, val.Statistics.LikeCount,
					val.Statistics.DislikeCount, t,
					"https://youtu.be/"+val.Id)
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

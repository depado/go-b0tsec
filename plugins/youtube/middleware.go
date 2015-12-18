package youtube

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"

	"google.golang.org/api/googleapi/transport"
	yt "google.golang.org/api/youtube/v3"
)

const (
	middlewareName = "youtube"
)

var ytre = regexp.MustCompile(`(?:https?://)?(?:(?:www\.)?youtube\.com/watch\?.*v=|youtu\.be/)([^&?]{11})`)

// Middleware is the youtube middleware.
type Middleware struct {
	Started bool
}

func init() {
	plugins.Middlewares[middlewareName] = new(Middleware)
}

// Get actually sends the data
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
	client := &http.Client{
		Transport: &transport.APIKey{Key: configuration.Config.GoogleAPIKey},
	}
	for _, bit := range strings.Fields(message) {
		rs := ytre.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			service, err := yt.New(client)
			if err != nil {
				fmt.Println("should send")
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

// Start starts the middleware and returns any occured error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		m.Started = true
	}
	return nil
}

// Stop stops the middleware and returns any occured error, nil otherwise
func (m *Middleware) Stop() error {
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

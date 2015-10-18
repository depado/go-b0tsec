package github

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const apiURL = "https://api.github.com/repos/%s/%s"

var re, _ = regexp.Compile("https?://github.com/([^/]+)/([^/]+)/?")

// Middleware is the github middleware
type Middleware struct{}

// RepoInfo represents the information relative to the repository
type RepoInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2])
			ri := RepoInfo{}
			err := utils.FetchURL(endpoint, &ri)
			if err == nil {
				ib.Privmsgf(configuration.Config.Channel, "%s : %v", ri.FullName, ri.Description)
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

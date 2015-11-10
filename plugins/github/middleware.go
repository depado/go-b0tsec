package github

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const apiURL = "https://api.github.com/repos/%s/%s%s"

var re = regexp.MustCompile("https?://github.com/([^/]+)/([^/]+)(/([^/]+)/[^/]+)?")

// Middleware is the github middleware
type Middleware struct{}

// RepoInfo represents the information relative to the repository
type RepoInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
}

type IssueInfo struct {
	Title string `json:"title"`
	State string `json:"state"`
}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			var extraStr string
			if len(rs[0][3]) > 0 {
				switch rs[0][4] {
				case "issues":
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], rs[0][3])
					ri := IssueInfo{}
					err := utils.FetchURL(endpoint, &ri)
					if err == nil {
						if ri.State == "open" {
							ri.State = "\x0303Opened\x03"
						} else {
							ri.State = "\x0304Closed\x0F\x03"
						}
						extraStr = fmt.Sprintf(" | %s - Status : %s", ri.Title, ri.State)
					}
				}
			}
			endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], "")
			ri := RepoInfo{}
			err := utils.FetchURL(endpoint, &ri)
			if err == nil {
				finalStr := fmt.Sprintf("%s : %v%s", ri.FullName, ri.Description, extraStr)
				ib.Privmsgf(configuration.Config.Channel, finalStr)
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

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

var re = regexp.MustCompile("https?://github.com/([^/]+)/([^/]+)(/(issues/[0-9]+|commit/[[:xdigit:]]{40}))?")

// Middleware is the github middleware
type Middleware struct{}

// RepoInfo represents the information relative to the repository
type RepoInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
}

type IssueInfo struct {
	Title  string `json:"title"`
	State  string `json:"state"`
	Number int    `json:"number"`
}

type CommitInfo struct {
	Commit struct {
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			var extraStr string
			if len(rs[0][3]) > 0 {
				if strings.HasPrefix(rs[0][4], "issues/") {
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], rs[0][3])
					ri := IssueInfo{}
					err := utils.FetchURL(endpoint, &ri)
					if err == nil {
						if ri.State == "open" {
							ri.State = "\x0303Opened\x03"
						} else {
							ri.State = "\x0304Closed\x0F\x03"
						}
						extraStr = fmt.Sprintf(" | #%4d %s - Status : %s", ri.Number, ri.Title, ri.State)
					}
				} else if strings.HasPrefix(rs[0][4], "commit/") {
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], strings.Replace(rs[0][3], "/commit/", "/commits/", 1))
					ri := CommitInfo{}
					err := utils.FetchURL(endpoint, &ri)
					if err == nil {
						extraStr = fmt.Sprintf(" | %s <%s> committed “%v”", ri.Commit.Author.Name, ri.Commit.Author.Email, strings.Replace(ri.Commit.Message, "\n", " ", -1))
					}
				}
			}
			endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], "")
			ri := RepoInfo{}
			err := utils.FetchURL(endpoint, &ri)
			if err == nil {
				var finalStr string
				if len(extraStr) > 0 {
					finalStr = fmt.Sprintf("%s %s", ri.FullName, extraStr)
				} else {
					finalStr = fmt.Sprintf("%s : %v", ri.FullName, ri.Description)
				}
				ib.Privmsgf(configuration.Config.Channel, finalStr)
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

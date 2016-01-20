package github

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/plugins"
	"github.com/depado/go-b0tsec/utils"
	"github.com/thoj/go-ircevent"
)

const (
	middlewareName = "github"
	apiURL         = "https://api.github.com/repos/%s/%s%s"
)

// Middleware is the github middleware
type Middleware struct {
	Started bool
}

var re = regexp.MustCompile("https?://github.com/([^/]+)/([^/]+)(/(issues/[0-9]+|commit/[[:xdigit:]]{40}|pull/[0-9]+))?")

func init() {
	plugins.Middlewares[middlewareName] = new(Middleware)
}

// Get actually sends the data
func (m *Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if !m.Started {
		return
	}
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			var extraStr string
			if len(rs[0][3]) > 0 {
				if strings.HasPrefix(rs[0][4], "issues/") {
					// If issue, get its info
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], rs[0][3])
					ri := IssueInfo{}
					if err := utils.FetchURL(endpoint, &ri); err == nil {
						status := ri.FormatStatus()
						extraStr = fmt.Sprintf("| #%4d %s - Status : %s", ri.Number, ri.Title, status)
					}
				} else if strings.HasPrefix(rs[0][4], "commit/") {
					// If commit, get its info
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], strings.Replace(rs[0][3], "/commit/", "/commits/", 1))
					ri := CommitInfo{}
					if err := utils.FetchURL(endpoint, &ri); err == nil {
						extraStr = fmt.Sprintf("| %s <%s> committed “%v”", ri.Commit.Author.Name, ri.Commit.Author.Email, strings.Replace(ri.Commit.Message, "\n", " ", -1))
					}
				} else if strings.HasPrefix(rs[0][4], "pull/") {
					// If pull request, get its info
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], strings.Replace(rs[0][3], "/pull/", "/pulls/", 1))
					ri := PullRequestInfo{}
					if err := utils.FetchURL(endpoint, &ri); err == nil && len(ri.State) > 0 {
						status := ri.FormatStatus()
						extraStr = fmt.Sprintf("| %s PR %s by %s : \x0303+%d\x03 \x0304-%d\x0F\x03 in %d commits across %d files", status, ri.Title, ri.User.Login, ri.Additions, ri.Deletions, ri.Commits, ri.ChangedFiles)
					}
				}
			}
			endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], "")
			ri := RepoInfo{}
			if err := utils.FetchURL(endpoint, &ri); err == nil {
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

// Start starts the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Start() error {
	if utils.StringInSlice(middlewareName, configuration.Config.Middlewares) {
		m.Started = true
	}
	return nil
}

// Stop stops the middleware and returns any occurred error, nil otherwise
func (m *Middleware) Stop() error {
	m.Started = false
	return nil
}

// IsStarted returns the state of the middleware
func (m *Middleware) IsStarted() bool {
	return m.Started
}

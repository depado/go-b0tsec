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

var re = regexp.MustCompile("https?://github.com/([^/]+)/([^/]+)(/(issues/[0-9]+|commit/[[:xdigit:]]{40}|pull/[0-9]+))?")

// Middleware is the github middleware
type Middleware struct{}

// RepoInfo represents the information relative to the repository
type RepoInfo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
}

// IssueInfo represents the information relative to the issue
type IssueInfo struct {
	Title  string `json:"title"`
	State  string `json:"state"`
	Number int    `json:"number"`
}

// CommitInfo represents the information relative to the commit
type CommitInfo struct {
	Commit struct {
		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Message string `json:"message"`
	} `json:"commit"`
}

// PullRequestInfo represents the information relative to the pull request
type PullRequestInfo struct {
	State string `json:"state"`
	Title string `json:"title"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
	Merged       bool `json:"merged"`
	Mergeable    bool `json:"mergeable"`
	Commits      int  `json:"commits"`
	Additions    int  `json:"additions"`
	Deletions    int  `json:"deletions"`
	ChangedFiles int  `json:"changed_files"`
}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			var extraStr string
			if len(rs[0][3]) > 0 {
				if strings.HasPrefix(rs[0][4], "issues/") {
					// If issue, get its info
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], rs[0][3])
					ri := IssueInfo{}
					err := utils.FetchURL(endpoint, &ri)
					if err == nil {
						if ri.State == "open" {
							ri.State = "\x0303Opened\x03"
						} else {
							ri.State = "\x0304Closed\x03"
						}
						extraStr = fmt.Sprintf(" | #%4d %s - Status : %s", ri.Number, ri.Title, ri.State)
					}
				} else if strings.HasPrefix(rs[0][4], "commit/") {
					// If commit, get its info
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], strings.Replace(rs[0][3], "/commit/", "/commits/", 1))
					ri := CommitInfo{}
					err := utils.FetchURL(endpoint, &ri)
					if err == nil {
						extraStr = fmt.Sprintf(" | %s <%s> committed “%v”", ri.Commit.Author.Name, ri.Commit.Author.Email, strings.Replace(ri.Commit.Message, "\n", " ", -1))
					}
				} else if strings.HasPrefix(rs[0][4], "pull/") {
					// If pull request, get its info
					endpoint := fmt.Sprintf(apiURL, rs[0][1], rs[0][2], strings.Replace(rs[0][3], "/pull/", "/pulls/", 1))
					ri := PullRequestInfo{}
					err := utils.FetchURL(endpoint, &ri)
					fmt.Println(ri)
					if err == nil && len(ri.State) > 0{
						var status string
						if ri.State == "open" {
							if ri.Mergeable {
								status = "\x0303Mergeable\x03"
							} else {
								status = "\x0304Unmergeable\x03"
							}
						} else {
							if ri.Merged {
								status = "\u0002Merged\u000F"
							} else {
								status = "\x0304Closed and not merged\x03"
							}
						}
						extraStr = fmt.Sprintf(" | %s PR %s by %s : \x0303+%d\x03 \x0304-%d\x03 in %d commits across %d files", status, ri.Title, ri.User.Login, ri.Additions, ri.Deletions, ri.Commits, ri.ChangedFiles)
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

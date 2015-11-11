package github

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

// FormatStatus returns a formatted (colored) string about the state of an IssueInfo
func (i IssueInfo) FormatStatus() string {
	if i.State == "open" {
		return "\x0303Opened\x03"
	}
	return "\x0304Closed\x0F\x03"
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

// FormatStatus returns a formatted (colored) string about the state of a PullRequestInfo
func (p PullRequestInfo) FormatStatus() string {
	if p.State == "open" {
		if p.Mergeable {
			return "\x0303Mergeable\x03"
		}
		return "\x0304Unmergeable\x0F\x03"
	}
	if p.Merged {
		return "\u0002Merged\u000F"
	}
	return "\x0304Closed and not merged\x0F\x03"
}

package title

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"

	"github.com/thoj/go-ircevent"
)

var re = regexp.MustCompile(`(?:https?://)(?:www.)?([^/]*).*`)

// Middleware is the github middleware
type Middleware struct{}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	cnf := configuration.Config
	if to == cnf.BotName {
		to = from
	}
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			host := rs[0][1]
			m := cnf.Middlewares
			// If middlewares youtube or github are disabled, we still get the
			// title of these sites.
			if (host == "youtube.be" || host == "youtube") && !utils.StringInSlice("youtube", m) ||
				host == "github.com" && !utils.StringInSlice("github", m) {
				resp, err := http.Get(rs[0][0])
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()
				fURL := resp.Request.URL.String()
				z := html.NewTokenizer(resp.Body)
				for {
					tt := z.Next()
					switch tt {
					case html.ErrorToken:
						return
					case html.StartTagToken:
						t := z.Token()
						if t.Data == "title" {
							tt = z.Next()
							t = z.Token()
							d := t.Data
							if len(d) > 450 {
								d = d[:450]
							}
							if fURL != rs[0][0] {
								ib.Privmsgf(to, "%v (%v)", d, fURL)
							} else {
								ib.Privmsg(to, d)
							}
							return
						}
					}
				}
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

package title

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/thoj/go-ircevent"
)

var re = regexp.MustCompile(`(?:https?://)(?:www.)?([^/]*).*`)

// Middleware is the github middleware
type Middleware struct{}

// Get actually sends the data
func (m Middleware) Get(ib *irc.Connection, from string, to string, message string) {
	if to == configuration.Config.BotName {
		to = from
	}
	for _, bit := range strings.Fields(message) {
		rs := re.FindAllStringSubmatch(bit, -1)
		if len(rs) > 0 {
			tld := rs[0][1]
			if tld != "youtube.com" && tld != "youtu.be" && tld != "github.com" {
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

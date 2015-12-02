package title

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"golang.org/x/net/html"

	"github.com/depado/go-b0tsec/configuration"
	"github.com/depado/go-b0tsec/utils"

	"github.com/rakyll/magicmime"
	"github.com/thoj/go-ircevent"
)

var re = regexp.MustCompile(`(?:https?://)(?:www.)?([^/]*).*`)

// Middleware is the github middleware
type Middleware struct{}

// Get the title token of a HTML page
func GetTitle(resp *http.Response, url string) string {
	fURL := resp.Request.URL.String()
	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return ""
		case html.StartTagToken:
			t := z.Token()
			if t.Data == "title" {
				tt = z.Next()
				t = z.Token()
				d := t.Data
				if len(d) > 450 {
					d = d[:450]
				}
				if fURL != url {
					return fmt.Sprintf("%v (%v)", d, fURL)
				} else {
					return d
				}
				return ""
			}
		}
	}
}

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
			if (host != "youtu.be" && host != "youtube.com" || !utils.StringInSlice("youtube", m)) &&
				(host != "github.com" || !utils.StringInSlice("github", m)) {

				// Avoid KeepAlives slowing down a simple GET, also we dont care
				// about unsafe TLS
				tr := &http.Transport{
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DisableKeepAlives: true,
				}
				client := &http.Client{Transport: tr}

				resp, err := client.Get(rs[0][0])
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()

				if err := magicmime.Open(magicmime.MAGIC_NO_CHECK_COMPRESS); err != nil {
					log.Println(err)
					return
				}
				defer magicmime.Close()

				data, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println(err)
					return
				}

				mimetype, err := magicmime.TypeByBuffer(data)
				if err != nil {
					log.Println(err)
					return
				}

				if strings.HasPrefix(mimetype, "HTML document") {
					resp.Body = ioutil.NopCloser(bytes.NewBuffer(data))
					if r := GetTitle(resp, rs[0][0]); r != "" {
						ib.Privmsg(to, r)
						return
					}
				}
				ib.Privmsgf(to, "File : (%s) %v", utils.HumanReadableSize(len(data)), mimetype)
			}
		}
	}
}

// NewMiddleware returns a new Middleware
func NewMiddleware() *Middleware {
	return new(Middleware)
}

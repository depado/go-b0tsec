package youtube

import (
	"fmt"
	"strings"

	yt "google.golang.org/api/youtube/v3"
)

type Middleware struct{}

type Plugin struct{}

func FormatOutput(v *yt.Video) string {
	t := strings.Replace(v.ContentDetails.Duration[2:len(v.ContentDetails.Duration)-1], "M", ":", -1)
	t = strings.Replace(t, "H", ":", -1)

	return fmt.Sprintf("%v [\x0303%v\x03 | \x0304%v\x0F\x03] (%v)",
		v.Snippet.Title, v.Statistics.LikeCount,
		v.Statistics.DislikeCount, t)
}

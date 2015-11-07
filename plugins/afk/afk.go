package afk

import "time"

// Data represents the time and the reason why a person if afk.
type Data struct {
	Since  time.Time
	Reason string
}

// Map contains the people being afk at the moment
var Map = make(map[string]Data)

package seen

import "time"

// Map contains the data about the last time the bot as seen someone
var Map = make(map[string]time.Time)

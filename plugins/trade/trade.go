package trade

// Cryptonator represents a structured response from the Cryptonator API
type Cryptonator struct {
	Ticker struct {
		Base    string `json:"base"`
		Target  string `json:"target"`
		Price   string `json:"price"`
		Volume  string `json:"volume"`
		Change  string `json:"change"`
		Markets []struct {
			Market string  `json:"market"`
			Price  string  `json:"price"`
			Volume float64 `json:"volume"`
		} `json:"markets"`
	} `json:"ticker"`
	Timestamp int    `json:"timestamp"`
	Success   bool   `json:"success"`
	Error     string `json:"error"`
}

package markovchains

import (
	"math/rand"
	"strings"
	"time"
)

// PrefixLen is the number of words per Prefix defined as the key for the map.
const PrefixLen = 2

// MainChain is the chain that will be available outside the package.
var MainChain *Chain

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	Key   string
	Chain map[string][]string
}

// Build builds the chain using the given string parameter
func (c *Chain) Build(s string) {
	p := make(Prefix, PrefixLen)
	for _, v := range strings.Split(s, " ") {
		key := p.String()
		c.Chain[key] = append(c.Chain[key], v)
		p.Shift(v)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate() string {
	p := make(Prefix, PrefixLen)
	var words []string
	for {
		choices := c.Chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[rand.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(key string) *Chain {
	return &Chain{key, make(map[string][]string)}
}

// Init Initializes the markov chain
func Init() error {
	var err error
	if err = InitBucketIfNotExists(); err != nil {
		return err
	}
	MainChain, err = GetChain("main")
	if err != nil {
		MainChain = NewChain("main")
	}
	rand.Seed(time.Now().UnixNano())
	return nil
}

package trade

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseOptions(t *testing.T) {
	payloads := map[string]options{
		"btc eth":                   {false, false, "btc", "eth", 1, ""},
		"2 btc eth":                 {false, false, "btc", "eth", 2, ""},
		"btc eth --market=kraken":   {false, false, "btc", "eth", 1, "kraken"},
		"2 btc eth --market=kraken": {false, false, "btc", "eth", 2, "kraken"},
		"eth btc --market=all":      {false, false, "eth", "btc", 1, "all"},
		"eth btc --sort=price":      {true, false, "eth", "btc", 1, ""},
	}

	for in, expected := range payloads {
		got, _ := parseOptions(strings.Split(in, " "))
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Expected %+v, got %+v", expected, got)
		}
	}
}

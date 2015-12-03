package utils

import "testing"

var teststr = []string{
	"The", "quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog.",
}

type ret struct {
	index int
	found bool
}

func TestIndexStringInSlice(t *testing.T) {
	payload := map[string]ret{
		"quick": ret{1, true},
		"brawn": ret{-1, false},
	}
	for key, val := range payload {
		i, f := IndexStringInSlice(key, teststr)
		if i != val.index || f != val.found {
			t.Errorf("For %v : Expected %v, got %v", key, val, ret{i, f})
		}
	}
}

func TestStringInSlice(t *testing.T) {
	payload := map[string]bool{
		"quick": true,
		"brawn": false,
		"dog.":  true,
		"jumps": true,
	}
	for key, val := range payload {
		f := StringInSlice(key, teststr)
		if f != val {
			t.Errorf("For %v : Expected %v, got %v", key, val, f)
		}
	}
}

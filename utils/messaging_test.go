package utils

import (
	"reflect"
	"testing"
)

func TestSplitMessage(t *testing.T) {
	var long string
	var short string
	var expected []string
	var parsed []string

	long = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas consectetur ipsum in urna ullamcorper, a lobortis mauris porta. Vivamus nec nibh nisi. Vestibulum consequat felis nulla. Pellentesque lectus sapien, lobortis quis finibus vitae, mollis a purus. Praesent aliquam orci sit amet ullamcorper egestas. Nam iaculis augue eu porttitor suscipit. Sed tempus massa nunc, id tincidunt leo aliquet vitae. In in ligula scelerisque, ornare sapien sed, rutrum ex. Ut fringilla commodo mauris, quis feugiat nunc tincidunt vitae. Ut sit amet ligula consectetur, varius nisi nec, aliquet metus. Aenean id urna lorem. Praesent non augue accumsan, tristique nunc nec, efficitur risus. Pellentesque tempus malesuada diam sit amet sollicitudin."
	expected = []string{"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas consectetur ipsum in urna ullamcorper, a lobortis mauris porta. Vivamus nec nibh nisi. Vestibulum consequat felis nulla. Pellentesque lectus sapien, lobortis quis finibus vitae, mollis a purus. Praesent aliquam orci sit amet ullamcorper egestas. Nam iaculis augue eu porttitor suscipit. Sed tempus massa nunc, id tincidunt leo aliquet vitae. In in ligula scelerisque, ornare sapien", "sed, rutrum ex. Ut fringilla commodo mauris, quis feugiat nunc tincidunt vitae. Ut sit amet ligula consectetur, varius nisi nec, aliquet metus. Aenean id urna lorem. Praesent non augue accumsan, tristique nunc nec, efficitur risus. Pellentesque tempus malesuada diam sit amet sollicitudin."}
	parsed = SplitMessage(long)
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
	short = "Toto !"
	expected = []string{"Toto !"}
	parsed = SplitMessage(short)
	if !reflect.DeepEqual(parsed, expected) {
		t.Errorf("Expected %v, got %v", expected, parsed)
	}
}

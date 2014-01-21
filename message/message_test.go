package message

import (
	"testing"
	"fmt"
)

var msg = Message{"woo", "test\n"}

func TestNameHashFunc(t *testing.T) {
	hash := hashName("test")
	hash2 := hashName("test")

	if hash != hash2 {
		t.Errorf("Hashes are not the same! %v %v", hash, hash2)
	}
}

func TestColor(t *testing.T) {
	c := getColor("lol")
	c2 := getColor("lol")

	if c != c2 {
		t.Errorf("color strings are not the same oh dear %v %v", c, c2)
	}

}

func TestFormatted(t *testing.T) {
	str := fmt.Sprintf("%s: %s\n", getPrefix("woo"), "test")

	if str != msg.Formatted() {
		t.Errorf("Formatted output is wrong! %v %s", str, msg.Formatted())
	}
}

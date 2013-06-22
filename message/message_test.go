package message

import (
	"testing"
)

var msg = Message{"woo", "test"}

func TestNameHashFunc(t *testing.T) {
	hash, _ := hashName("test")
	hash2, _ := hashName("test")

	if hash != hash2 {
		t.Errorf("Hashes are not the same! %v %v", hash, hash2)
	}
}

func TestColor(t *testing.T) {

}

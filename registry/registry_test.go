package registry

import (
	"testing"
)


func TestReadAndWrite(t *testing.T) {
	reg := New()

	reg.Set("foo", 1)

	c := reg.Get("foo")
	if c != 1 {
		t.Errorf("Expected %v!", c)
	}
}

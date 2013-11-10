package main

import (
	"testing"
)

func TestReadAndWrite(t *testing.T) {
	reg := NewRegistry()

	reg.Set("foo", 1)

	c, _ := reg.Get("foo")
	if c != 1 {
		t.Errorf("Expected %v!", c)
	}
}

func TestGetStatus(t *testing.T) {
	reg := NewRegistry()

	_, status := reg.Get("foo")

	if status {
		t.Errorf("Expected status to be false!")
	}
}

func TestAnyType(t *testing.T) {
	reg := NewRegistry()
	reg.Set("foo", 1)

	reg2 := NewRegistry()
	reg2.Set("foo", "bar")

	v1, _ := reg.Get("foo")
	v2, _ := reg2.Get("foo")
	if v1 != 1 {
		t.Error("Expected value to be 1")
	}

	if v2 != "bar" {
		t.Error("Expected value to be 'bar'")
	}

}

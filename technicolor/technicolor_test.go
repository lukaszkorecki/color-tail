package technicolor

import (
	"fmt"
	"testing"
)

func TestRedPaint(t *testing.T) {

	s := "test"
	exp := "\033[31mtest\033[39m"
	ns := Paint(s, "red")

	if ns != exp {
		t.Errorf("expected %v to equal %v", ns, exp)
	} else {
		fmt.Printf("%v | %v", ns, exp)
	}
}

func TestRandPaint(t *testing.T) {
	s := "test"
	strings := make([]string, len(Colors))
	for _, color := range Colors {
		strings = append(strings, Paint(s, color))
	}

	sr := RandPaint(s)
	hasString := false

	// search through all possible colored strings
	// and check if RandPaint generated one of them
	for _, str := range strings {
		if !hasString {
			hasString = (sr == str)
		}
	}
	if !hasString {
		t.Errorf("RandPaint didn't generate valid colored string: %v", sr)
	}
}

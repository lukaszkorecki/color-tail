package message

import (
	"../technicolor"
	"../registry"
	"fmt"
	"strings"
	"io"
	"crypto/sha1"
)

type Message struct {
	Name  string
	Event string
}

var (
	h = sha1.New()
	nameMap = registry.New()
	colorMap = registry.New()
)


// hashes name once and stores it in name map.
// which is suitable for consumption
func hashName(name string) string {
	hash, ok := nameMap.Get(name)
	if ! ok {
		io.WriteString(h, name)
		v := fmt.Sprintf("%x", h.Sum(nil)[0:3])
		nameMap.Set(name, v)
		hash = v
	}

	// force string casting of AnyType
	return hash.(string)
}

func getColor(name string) string {
	color, ok := colorMap.Get(name)
	if ! ok {
		color = technicolor.RandColorName()
		colorMap.Set(name, color)
	}
	return color.(string)
}


func (m Message) Print() {
	hn := hashName(m.Name)
	color := getColor(m.Name)

	for _, s := range strings.Split(m.Event, "\n") {
		n := technicolor.Paint(hn, color)
		fmt.Printf("%v: %v\n", n, s)
	}
}

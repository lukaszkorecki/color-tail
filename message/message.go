package message

import (
	"crypto/sha1"
	"fmt"
	r "github.com/lukaszkorecki/color-tail/registry"
	t "github.com/lukaszkorecki/color-tail/technicolor"
	"io"
	"strings"
)

type Message struct {
	Name  string
	Event string
}

var (
	h        = sha1.New()
	nameMap  = r.New()
	colorMap = r.New()
)

// hashes name once and stores it in name map.
// so for /var/log/nginx/error.log it will return 5faabc4e
func hashName(name string) string {
	hash, ok := nameMap.Get(name)
	if !ok {
		io.WriteString(h, name)
		v := fmt.Sprintf("%x", h.Sum(nil)[0:3])
		nameMap.Set(name, v)
		hash = v
	}

	// force string casting of AnyType
	return hash.(string)
}

// Pick a random color once and then store it for given file name
func getColor(name string) string {
	color, ok := colorMap.Get(name)
	if !ok {
		color = t.RandColorName()
		colorMap.Set(name, color)
	}
	return color.(string)
}

func getPrefix(name string) string {
	hn := hashName(name)
	color := getColor(name)
	prefix := t.Paint(hn, color)
	return prefix
}

func formatEvent(prefix, event string) string {
	lines := strings.Split(event, "\n")
	length := len(lines) // last element is \n

	buf := make([]string, length)

	for i, line := range lines {
		s := fmt.Sprintf("%v: %v\n", prefix, line)
		if i < length-1 {
			buf = append(buf, s)
		}
	}

	return strings.Join(buf, "")
}

// Print a file message and color code the file name hash
func (m Message) Print() {
	prefix := getPrefix(m.Name)
	str := formatEvent(prefix, m.Event)

	fmt.Print(str)
}

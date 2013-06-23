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
// so for /var/log/nginx/error.log it will return 5faabc4e
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

// Pick a random color once and then store it for given file name
func getColor(name string) string {
	color, ok := colorMap.Get(name)
	if ! ok {
		color = technicolor.RandColorName()
		colorMap.Set(name, color)
	}
	return color.(string)
}

func getPrefix(name string) string {
	hn := hashName(name)
	color := getColor(name)
	prefix := technicolor.Paint(hn, color)
	return prefix
}

func formatEvent(prefix, event string) string {
	lines := strings.Split(event, "\n")

	buf := make([]string, len(lines))

	for _, line := range lines {
		s := fmt.Sprintf("%v: %v", prefix, line)
		fmt.Printf("%v -%v-", line, line)
		if len(s) > 0 {
			buf = append(buf, s)
		}
	}

	return strings.Join(buf, "\n")
}

// Print a file message and color code the file name hash
func (m Message) Print() {
	prefix := getPrefix(m.Name)
	str := formatEvent(prefix, m.Event)

	fmt.Print(str)
}

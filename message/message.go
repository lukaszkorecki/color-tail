package message

import (
	"../technicolor"
	"fmt"
	"strings"
	"io"
	"crypto/md5"
)

type Message struct {
	Name  string
	Event string
}

var h = md5.New()

func hashName(name string) string {
	io.WriteString(h, name)
	return fmt.Sprintf("%x", h.Sum(nil)[0:3])
}

func (m Message) Print() {
	for _, s := range strings.Split(m.Event, "\n") {
		hn := hashName(m.Name)
		n := technicolor.RandPaint(hn)
		fmt.Printf("%v: %v\n", n, s)
	}
}

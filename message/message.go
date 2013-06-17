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

func hashId(name string) int {
	io.WriteString(h, name)
	return int(h.Sum(nil)[0])
}

func hashName(name string) string {
	io.WriteString(h, name)
	return fmt.Sprintf("%x", h.Sum(nil)[0:3])
}


func (m Message) Print() {
	hn := hashName(m.Name)
	n := technicolor.RandPaint(hn)
	fmt.Printf("%v:\n", technicolor.RandPaint(m.Name))

	for _, s := range strings.Split(m.Event, "\n") {
		if len(s) > 0 {
			fmt.Printf("%v: %v\n", n, s)
		}
	}
}

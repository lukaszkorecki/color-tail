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

func hashName(name string) (string, string){
	io.WriteString(h, name)
	s := fmt.Sprintf("%x", h.Sum(nil)[0:3])
	return s, technicolor.RandColorName()
}


func (m Message) Print() {
	hn, color := hashName(m.Name)
	for _, s := range strings.Split(m.Event, "\n") {
		n := technicolor.Paint(hn, color)
		fmt.Printf("%v: %v\n", n, s)
	}
}

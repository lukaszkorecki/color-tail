package message

import (
	"../technicolor"
	"../registry"
	"fmt"
	"strings"
	"io"
	"crypto/md5"
)

type Message struct {
	Name  string
	Event string
}

var (
	h = md5.New()
	nameMap = registry.New()
)


// hashes name once and stores it in name map.
// TODO move this to a separate module
func hashName(name string) (string, string){
	hash, ok := nameMap.Get(name)
	if ! ok {
		io.WriteString(h, name)
		v := fmt.Sprintf("%x", h.Sum(nil)[0:3])
		nameMap.Set(name, v)
		hash = v
	}
	// force string casting of AnyType
	r := fmt.Sprintf("%v", hash)
	return r, technicolor.RandColorName()
}


func (m Message) Print() {
	hn, color := hashName(m.Name)
	for _, s := range strings.Split(m.Event, "\n") {
		n := technicolor.Paint(hn, color)
		fmt.Printf("%v: %v\n", n, s)
	}
}

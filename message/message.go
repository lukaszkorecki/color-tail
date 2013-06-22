package message

import (
	"../technicolor"
	"fmt"
	"strings"
	"io"
	"crypto/md5"
	"sync"
)

type Message struct {
	Name  string
	Event string
}

type nameHashMap struct {
	lock sync.RWMutex
	store map[string]string
}
func (n *nameHashMap) Get(key string) (string, bool) {
	n.lock.RLock()
	defer n.lock.RUnlock()
	v, ok := n.store[key]
	return v, ok
}
func (n *nameHashMap) Set(key, val string) string {
	n.lock.Lock()
	defer n.lock.Unlock()
	n.store[key] = val
	return val
}

var (
	h = md5.New()
	nameMap = &nameHashMap{store: make(map[string]string)}
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
	return hash, technicolor.RandColorName()
}


func (m Message) Print() {
	hn, color := hashName(m.Name)
	for _, s := range strings.Split(m.Event, "\n") {
		n := technicolor.Paint(hn, color)
		fmt.Printf("%v: %v\n", n, s)
	}
}

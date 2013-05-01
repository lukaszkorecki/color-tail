package message

import (
	"fmt"
	"log"
)

type Message struct {
	Name  string
	Event string
}

func (m Message) Print() {
	log.Printf(fmt.Sprintf("%v: %v", m.Event, m.Name))
}

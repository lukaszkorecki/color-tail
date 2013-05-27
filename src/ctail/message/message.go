package message

import (
	"fmt"
)

type Message struct {
	Name  string
	Event string
}

func (m Message) Print() {
	fmt.Printf("%v:\n%v\n", m.Name, m.Event)
}

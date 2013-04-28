package message
import (
  "log"
  "fmt"
)

type Message  struct {
  Name string
  Event string
}

func (m Message) Print() {
  log.Printf(fmt.Sprintf("%v: %v", m.Event, m.Name))
}


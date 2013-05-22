package main

import (
	"ctail/message"
	"github.com/howeyc/fsnotify"
	"fmt"
	"log"
	"os"
  "bytes"
	"path/filepath"
)

// this map holds references to last position for given file
// and it needs to be updated whenever a file is read...
// XXX what about two goroutines updating this map at the same time?
// we need to take this blog post into account:
// - http://blog.golang.org/2013/02/go-maps-in-action.html
// especially the section about maps not being thread safe (which is easy to fix)
var Registry = make(map[string]int64)


func colorize(stuff string) string {
  colorTable := make(map[string]string)

  colorTable["black"]   = "\033[30m"
  colorTable["red"]     = "\033[31m"
  colorTable["green"]   = "\033[32m"
  colorTable["yellow"]  = "\033[33m"
  colorTable["blue"]    = "\033[34m"
  colorTable["magenta"] = "\033[35m"
  colorTable["cyan"]    = "\033[36m"
  colorTable["white"]   = "\033[37m"
  colorTable["reset"]   = "\033[39m"

  return colorTable["red"] + stuff + colorTable["reset"]

}
func fileNotification(fname string) message.Message {
  file, err := os.Open(fname)

  if err != nil {
    return message.Message{fname, "Can't open file!"}
  }

	lastPosition := Registry[fname]
  buf := bytes.NewBuffer(make ([]byte, lastPosition))

  _ , err = buf.ReadFrom(file)
  if err != nil {
    log.Printf("Reading from %v failed: %v", fname, err)
  }

  length := int64(buf.Len())

  offset :=  length - lastPosition

  if offset == length { offset = 8 }

  log.Printf("length: %v lastPosition: %v offset: %v",length, lastPosition, offset)

  str := string(buf.Next(int(offset)))
  Registry[fname] = length
	return message.Message{colorize(fname), str}
}


func monitorPath(fname string, notify chan message.Message) {
	watcher, _ := fsnotify.NewWatcher()
	notify <- message.Message{fname, "Start!"}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				log.Printf("<<<<< %v", event.Name)
				notify <- fileNotification(event.Name)
			case err := <-watcher.Error:
				notify <- message.Message{fname, fmt.Sprintf("Error: %v", err)}
				watcher.Close()
			}
		}
	}()

	watcher.Watch(fname)
}

func main() {
	out := make(chan message.Message)
	if len(os.Args) == 1 {
		log.Fatal("Needs files!")
		os.Exit(1)
	}

	for i := 1; i < len(os.Args); i++ {
		fname, _ := filepath.Abs(os.Args[i])
		// XXX check if file is readable?
		Registry[fname] = 0
		go monitorPath(fname, out)
	}

	for {
		message := <-out
		message.Print()
	}
}

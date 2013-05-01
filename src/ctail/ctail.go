package main

import (
	"ctail/message"
	"fmt"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// this map holds references to last position for given file
// and it needs to be updated whenever a file is read...
// XXX what about two goroutines updating this map at the same time?
// we need to take this blog post into account:
// - http://blog.golang.org/2013/02/go-maps-in-action.html
// especially the section about maps not being thread safe (which is easy to fix)
var Registry = make(map[string]int)

func fileNotification(fname string) message.Message {
	lastPosition := Registry[fname]
	Registry[fname] = lastPosition + 50
	contents, _ := ioutil.ReadFile(fname)
	return message.Message{fname, string(contents)}
}

func monitorPath(fname string, notify chan message.Message) {
	watcher, _ := fsnotify.NewWatcher()
	notify <- message.Message{fname, "Start!"}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				log.Printf(event.Name)
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

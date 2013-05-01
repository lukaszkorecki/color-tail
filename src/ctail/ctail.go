package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"log"
	"message"
	"os"
	"path/filepath"
)

func fileNotification(fname string) message.Message {
	contents, _ := ioutil.ReadFile(fname)
	s := fmt.Sprintf("Modified! %v", contents)
	return message.Message{fname, s}
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
		log.Printf(fname)
		go monitorPath(fname, out)
	}

	for {
		message := <-out
		message.Print()
	}
}

package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
  "io/ioutil"
	"path/filepath"
)

func fileNotification(fname string) string {
	s := fmt.Sprintf("Modified! %v", fname)
  contents, _ := ioutil.ReadFile(fname)
  log.Printf(fmt.Sprintf("%v", contents))
	return s
}

func monitorPath(fname string, notify chan string) {
	watcher, _ := fsnotify.NewWatcher()
	notify <- "Watching " + fname + "!"

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				log.Printf(event.Name)
        notify <- fileNotification(event.Name)
			case err := <-watcher.Error:
				notify <- fmt.Sprintf("%v", err)
        watcher.Close()
			}
		}
	}()

	watcher.Watch(fname)
}

func main() {
	out := make(chan string)
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
		log.Println(message)
	}
}

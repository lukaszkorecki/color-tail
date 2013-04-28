package main

import (
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

func monitorPath(fname string, notify chan string) {
	log.Printf("Watching %v", fname)
	watcher, _ := fsnotify.NewWatcher()
	notify <- "Ready!"

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				notify <- event.Name
			case err := <-watcher.Error:
				log.Printf("ERROR %v", err)
				notify <- "error!"
				watcher.Close()
			}
		}
	}()

	watcher.Watch(fname)
}

func main() {
	out := make(chan string)

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

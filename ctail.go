package main

import (
	"./message"
	"./filemonitor"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// keeps a track of file size...
// XXX is it good to have it as a global?
var (
	version = "0.1.0"
)

// the meat of the programme
// whenever an event of a file change is received we check it's previous
// size (if we have it) and then extract the lines added and pack them
// into a Message to print.
func fileChanged(fname string) message.Message {
	return filemonitor.Changed(fname)
}

// main... event handler so to speak
func monitorPath(fname string, notify chan message.Message) {
	watcher, _ := fsnotify.NewWatcher()
	watcher.Watch(fname)

	log.Printf("Monitoring %v", fname)

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				notify <- fileChanged(event.Name)
			case err := <-watcher.Error:
				notify <- message.Message{fname, fmt.Sprintf("Error: %v", err)}
				watcher.Close()
			}
		}
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if len(os.Args) <= 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("ctail v.%v\nUsage: ctail <path to files>\n", version)
		os.Exit(1)
	}

	out := make(chan message.Message)

	for _, argPath := range os.Args[1:] {
		paths, err := filepath.Glob(argPath)
		if err != nil {
			log.Fatalf("Invalid path! %v", argPath);
			os.Exit(1)
		}

		for _, path := range paths {
			fname, _ := filepath.Abs(path)

			if filemonitor.InitialSize(fname) {
				go monitorPath(fname, out)
			} else {
				log.Printf("!! File can't be read!: %v", fname)
				os.Exit(1)
			}
		}
	}

	for {
		message := <-out
		message.Print()
	}
}

package main

import (
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

// main... event handler so to speak
func monitorPath(fname string, notify chan Message) {
	watcher, _ := fsnotify.NewWatcher()
	watcher.Watch(fname)

	log.Printf("Monitoring %v", fname)

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				notify <- Changed(event.Name)
			case err := <-watcher.Error:
				notify <- Message{fname, fmt.Sprintf("Error: %v", err)}
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

	var filePaths []string

	// get file list from arguments, ignoring 0th arg since it's the
	// bin name
	for _, argPath := range os.Args[1:] {
		paths, err := filepath.Glob(argPath)
		if err != nil {
			log.Fatalf("Invalid path! %v", argPath)
			os.Exit(1)
		}

		for _, path := range paths {
			fname, _ := filepath.Abs(path)
			filePaths = append(filePaths, fname)
		}
	}

	cnt := 0
	out := make(chan Message)

	for _, fname := range filePaths {
		if InitialSize(fname) {
			go monitorPath(fname, out)
			cnt += 1
		} else {
			cnt -= 1
			log.Printf("!! File can't be read!: %v", fname)
		}
	}

	if cnt > 0 {
		for {
			message := <-out
			message.Print()
		}
	} else {
		log.Print("No files to monitor!")

		os.Exit(1)
	}
}

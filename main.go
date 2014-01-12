package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
	"runtime"
	fm "./filemonitor"
	m "./message"
)

var (
	version = "0.1.0a" // TODO read this from VERSION when compiling?
)

// main... event handler so to speak
func monitorPath(fname string, notify chan m.Message) {
	watcher, _ := fsnotify.NewWatcher()
	watcher.Watch(fname)

	log.Printf("Monitoring %v", fname)

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				notify <- fm.Changed(event.Name)
			case err := <-watcher.Error:
				notify <- m.Message{fname, fmt.Sprintf("Error: %v", err)}
				watcher.Close()
			}
		}
	}()
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// TODO FIXME use something sane for option parsing
	if len(os.Args) <= 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		name := os.Args[0]
		fmt.Printf(`%s v.%v

Usage: %s <path to files>
`, name, version, name)
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

	successCount := 0
	out := make(chan m.Message)

	for _, fname := range filePaths {
		if fm.InitialSize(fname) {
			go monitorPath(fname, out)
			successCount += 1
		} else {
			successCount -= 1
			log.Printf("!! File can't be read!: %v", fname)
		}
	}

	if successCount > 0 {
		for {
			message := <-out
			message.Print()
		}
	} else {
		log.Print("No files to monitor!")
		os.Exit(1)
	}
}

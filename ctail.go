package main

import (
	"./message"
	"./registry"
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// keeps a track of file size...
// XXX is it good to have it as a global?
var (
	sizeMap = registry.New()
	version = "0.1.0"
)

// FIXME crate a separate package for file ops?
func setInitialSize(fname string) bool {
	file, err := os.Open(fname)
	defer file.Close()

	if err != nil {
		log.Printf("!!! Can't open file: %v", fname)
		return false
	}


	size, statErr := getFileSize(file)

	if !statErr {
		log.Printf("!!! Can't file size!", fname)
		return false // file can't be read...
	}
	sizeMap.Set(fname, size)

	return true
}

func getFileSize(f *os.File) (int64, bool) {
	stat, err := f.Stat()
	if err != nil {
		return 0, false
	}

	return int64(stat.Size()), true
}

// the meat of the programme
// whenever an event of a file change is received we check it's previous
// size (if we have it) and then extract the lines added and pack them
// into a Message to print.
func fileChanged(fname string) message.Message {
	file, err := os.Open(fname)
	defer file.Close()

	// get file size
	size, statErr := getFileSize(file)

	if err != nil || statErr != true {
		return message.Message{fname, "Can't open file!"}
	}

	lastSize, _ := sizeMap.Get(fname)
	offset, _ := lastSize.(int64)

	// file got trimmed - or something reported wrong size
	if offset >= size || offset <= 0 {
		offset = int64(float64(size) / 0.25)
	}

	buf := make([]byte, offset)

	// read only recently appended content
	_, readErr := file.ReadAt(buf, offset)
	if readErr != nil && readErr != io.EOF {
		log.Printf("!!! Reading from %v failed: %v", fname, readErr)
	}

	// update file's size in the registry
	sizeMap.Set(fname, int64(size))
	return message.Message{fname, string(buf)}
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

			if setInitialSize(fname) {
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

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

// keeps a track of file size
var sizeMap = registry.New()

// the meat of the programme
// whenever an event of a file change is received we check it's previous
// size (if we have it) and then extract the lines added and pack them
// into a Message to print.
// If size was not detected we print last 10 lines <- TODO FIXME that's a lie
func fileChanged(fname string) message.Message {
	file, err := os.Open(fname)
	defer file.Close()

	// get file size
	stat, statErr := file.Stat()
	size := int64(stat.Size())

	if err != nil || statErr != nil {
		return message.Message{fname, "Can't open file!"}
	}

	lastSize, isSet := sizeMap.Get(fname)

	// assume we don't have size yet...
	offset := int64(0)
	if isSet {
		// Type assert here - we know that if it's set, it's int
		offset, _ = lastSize.(int64)
	}

	// file got trimmed - or something reported wrong size
	if offset >= size || offset <= 0 {
		offset = size / 2
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
	notify <- message.Message{fname, fmt.Sprintf("Start! %v", fname)}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				notify <- fileChanged(event.Name)
			case err := <-watcher.Error:
				notify <- message.Message{fname, fmt.Sprintf("Error: %v", err)}
				watcher.Close()
				// XXX FIXME remove file from watchers collection?
			}
		}
	}()

	watcher.Watch(fname)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	out := make(chan message.Message)

	if len(os.Args) == 1 {
		log.Fatal("Needs files!")
		os.Exit(1)
	}

	for i := 1; i < len(os.Args); i++ {
		fname, _ := filepath.Abs(os.Args[i])
		// XXX check if file is readable?
		// read file size here and include it in sizeMap?
		sizeMap.Set(fname, 0)
		go monitorPath(fname, out)
	}

	for {
		message := <-out
		message.Print()
	}
}

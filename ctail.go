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
var sizeMap = registry.New()

func setInitialSize(fname string) bool {
	file, err := os.Open(fname)
	defer file.Close()

	size, statErr := getFileSize(file)

	if err != nil || statErr != true {
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

		if setInitialSize(fname) {
			go monitorPath(fname, out)
		} else {
			log.Printf("!! File can't be read!: %v", fname)
		}
	}

	for {
		message := <-out
		message.Print()
	}
}

package main

import (
	"./message"
	"./registry"
	"fmt"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
)

var reg = registry.New()

func colorize(stuff string) string {
	colorTable := make(map[string]string)

	colorTable["black"] = "\033[30m"
	colorTable["red"] = "\033[31m"
	colorTable["green"] = "\033[32m"
	colorTable["yellow"] = "\033[33m"
	colorTable["blue"] = "\033[34m"
	colorTable["magenta"] = "\033[35m"
	colorTable["cyan"] = "\033[36m"
	colorTable["white"] = "\033[37m"
	colorTable["reset"] = "\033[39m"

	return colorTable["red"] + stuff + colorTable["reset"]

}
func fileChanged(fname string) message.Message {
	file, err := os.Open(fname)
	defer file.Close()

	// get file size
	stat, err2 := file.Stat()
	size := int64(stat.Size())

	if err != nil || err2 != nil {
		return message.Message{fname, "Can't open file!"}
	}

	lastSize := reg.Get(fname)
	offset := int64(0)

	if lastSize != 0 {
		offset = lastSize
	}

	// file got trimmed - or something reported wrong size
	if offset >= size || offset <= 0 {
		offset = size/2
	}


	log.Printf("lastSize: %v size: %v offset: %v", lastSize, size, offset)
	buf := make([]byte, offset)

	_, readErr := file.ReadAt(buf, offset)
	if readErr != nil {
		log.Printf("!!! Reading from %v failed: %v", fname, readErr)
	}

	reg.Set(fname, int64(size))
	return message.Message{colorize(fname), string(buf)}
}

func monitorPath(fname string, notify chan message.Message) {
	watcher, _ := fsnotify.NewWatcher()
	notify <- message.Message{fname, "Start!"}

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				log.Printf("<<<<< %v", event.Name)
				notify <- fileChanged(event.Name)
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
		reg.Set(fname, 0)
		go monitorPath(fname, out)
	}

	for {
		message := <-out
		message.Print()
	}
}
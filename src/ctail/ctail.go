package main
import (
  "log"
  "os"
  "github.com/howeyc/fsnotify"
)

func monitor(path string, watcher *fsnotify.Watcher) {
  go func() {
    for {
      select {
      case event := <-watcher.Event:
        log.Println("event: ", event)
      case err := <-watcher.Error:
        log.Println("error: ", err)
        log.Println("closing!", path)
        watcher.Close()

      }
    }
  }()

  watcher.Watch(path)
}

func main() {
  fname := ""
  if len(os.Args) == 2 {
    fname = os.Args[1]
  }
  log.Printf("Watching %v", fname)
  watcher, _ := fsnotify.NewWatcher()
  monitor(fname, watcher)
  for {
    log.Println("!")
  }

}

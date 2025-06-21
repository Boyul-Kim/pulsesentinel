package agent

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func Watch(paths []string) {
	if len(paths) < 1 {
		log.Fatal("must specify at least one path to watch")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("creating a new watcher: %s", err)
	}
	defer w.Close()

	go watchLoop(w)

	for _, p := range paths {
		err = w.Add(p)
		if err != nil {
			log.Fatalf("%q: %s", p, err)
		}
	}

	<-make(chan struct{}) // Block forever
}

func watchLoop(w *fsnotify.Watcher) {
	i := 0
	for {
		select {
		case err, ok := <-w.Errors:
			if !ok {
				return
			}
			log.Fatalf("ERROR: %s", err)
		// Read from Events.
		case e, ok := <-w.Events:
			if !ok {
				return
			}
			i++
			log.Printf("%3d %s %s", i, e)

			//eventually need to set up gRPC to send to event-ingestor service via client.go
		}
	}
}

package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/exp/inotify"
)

var watcher *inotify.Watcher

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func newWatcher() error {
	var err error
	watcher, err = inotify.NewWatcher()
	if err != nil {
		return err
	}
	err = watcher.AddWatch("/etc/ssl/server.key", inotify.IN_IGNORED)
	if err != nil {
		return err
	}
	return nil
}

func refreshWatcher() error {
	var err error
	err = watcher.Close()
	if err != nil {
		return err
	}
	return newWatcher()
}

func main() {
	err := newWatcher()
	if err != nil {
		log.Fatal(err)
	}

	errChan := make(chan error, 1)

	http.HandleFunc("/", HelloServer)

	go func() {
		errChan <- http.ListenAndServeTLS(":443", "/etc/ssl/server.crt", "/etc/ssl/server.key", nil)
	}()

	for {
		select {

		case err := <-errChan:
			if err != nil {
				log.Println("ListenAndServe: ", err)
			}
		case event := <-watcher.Event:
			log.Printf("event: %v", event)
			err := refreshWatcher()
			if err != nil {
				log.Println(err)
			}
		case err := <-watcher.Error:
			log.Println("error:", err)
		}

	}

}

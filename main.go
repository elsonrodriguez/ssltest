package main

import (
	"io"
	"log"
	"net/http"

	"golang.org/x/exp/inotify"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func main() {
	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.AddWatch("/etc/ssl/server.key", inotify.IN_MODIFY|inotify.IN_DELETE)
	if err != nil {
		log.Fatal(err)
	}
	err = watcher.AddWatch("/etc/ssl/server.crt", inotify.IN_MODIFY|inotify.IN_DELETE)
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
				log.Fatal("ListenAndServe: ", err)
			}
		case event := <-watcher.Event:
			log.Println("event:", event)
			
			err = watcher.AddWatch("/etc/ssl/server.key", inotify.IN_MODIFY|inotify.IN_DELETE)
			if err != nil {
				log.Fatal(err)
			}

		case err := <-watcher.Error:
			log.Println("error:", err)
		}

	}

}

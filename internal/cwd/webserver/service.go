package main

import (
	"context"
	"ela/foundation/path"
	"log"
	"net/http"
	"os"
)

const PORT = "80"

type MyService struct {
	server     *http.Server
	running    bool
	fileServer http.Handler
}

func (s *MyService) OnStart() error {
	s.running = true
	s.server = &http.Server{Addr: ":" + PORT}
	wwwPath := path.GetSystemWWW()
	s.fileServer = http.FileServer(http.Dir(wwwPath))
	indexFile := wwwPath + "/index.html"
	// handle any requests
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		// does the file exist? then serve the file
		if _, err := os.Stat(wwwPath + url); err == nil {
			s.fileServer.ServeHTTP(rw, r)
		} else { // hence use the index file
			http.ServeFile(rw, r, indexFile)
		}
	})
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal("Server error", err.Error())
		}
	}()
	return nil
}

func (s *MyService) OnEnd() error {
	return s.server.Shutdown(context.TODO())
}

func (s MyService) IsRunning() bool {
	return s.running
}

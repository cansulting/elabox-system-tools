package main

import (
	"context"
	"ela/foundation/path"
	"log"
	"net/http"
)

const PORT = "80"

type MyService struct {
	server  *http.Server
	running bool
}

func (s *MyService) OnStart() error {
	s.running = true
	s.server = &http.Server{Addr: ":" + PORT}
	wwwPath := path.GetSystemWWW()
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		http.FileServer(http.Dir(wwwPath)).ServeHTTP(rw, r)
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

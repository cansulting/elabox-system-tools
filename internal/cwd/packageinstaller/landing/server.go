package landing

import (
	"context"
	"log"
	"net/http"
	"time"
)

const PORT = "80"

var server *http.Server

const TIMEOUT = 10 // timeout for server initialization
var connected = 0

func Initialize(landingPagePath string) error {
	log.Println("Listening and serve @port", PORT, "www dir =", landingPagePath)

	server = &http.Server{Addr: ":" + PORT}
	// step: handle request
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		connected++
		http.FileServer(http.Dir(landingPagePath)).ServeHTTP(rw, r)
	})
	// step: listen and serve to target port
	go func() {
		elapsed := time.Now().Unix()
		for {
			err := server.ListenAndServe()
			if err == nil {
				break
			} else {
				log.Println("Issue found, retrying...", err.Error())
			}
			// step: check if this is waiting for too long
			diff := time.Now().Unix() - elapsed
			if diff > TIMEOUT {
				log.Fatal("Server error", err.Error())
				break
			}
			// sleep for a while
			time.Sleep(time.Millisecond * 500)
		}
	}()
	return nil
}

// wait for any users to connect to landing page
func WaitForConnection() {
	for connected == 0 {
		log.Println("Waiting @ port", PORT)
		time.Sleep(time.Second)
	}
	log.Println("Resuming...")
}

func Shutdown() error {
	return server.Shutdown(context.TODO())
}

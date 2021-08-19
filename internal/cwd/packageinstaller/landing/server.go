package landing

import (
	"ela/foundation/errors"
	"ela/foundation/system"
	"ela/server"
	"log"
	"net/http"
	"os"
	"time"
)

const PORT = "80"

var serverhandler *server.Manager

const TIMEOUT = 10 // timeout for server initialization
var connected = 0

func Initialize(landingPagePath string) error {
	// step: serve event server
	serverhandler = &server.Manager{}
	serverhandler.Setup()

	// step: init web server
	fileserver := http.FileServer(http.Dir(landingPagePath))
	log.Println("Landing page path =", landingPagePath)
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		connected++
		url := r.URL.Path
		if _, err := os.Stat(landingPagePath + url); err == nil {
			fileserver.ServeHTTP(rw, r)
		} else {
			http.Redirect(rw, r, "/", http.StatusFound)
		}
	})
	serverhandler.ListenAndServe()
	serverhandler.EventServer.SetStatus(system.UPDATING, nil)
	return nil
}

func GetServer() *server.Manager {
	return serverhandler
}

// wait for any users to connect to landing page
func WaitForConnection() {
	for connected == 0 {
		log.Println("Waiting @ port", PORT)
		time.Sleep(time.Second)
	}
	log.Println("Resuming...")
}

// use to shutdown the server
func Shutdown() error {
	log.Println("Shutting down event and server...")
	// close event server
	if serverhandler != nil {
		if err := serverhandler.Stop(); err != nil {
			return errors.SystemNew("Error closing connector.", err)
		}
	}
	return nil
}

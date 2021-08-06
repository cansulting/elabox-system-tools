package landing

import (
	"context"
	"ela/foundation/errors"
	"ela/foundation/event"
	"ela/foundation/system"
	"ela/internal/cwd/global/server"
	"ela/internal/cwd/system/global"
	"log"
	"net/http"
	"time"
)

const PORT = "80"

var webserver *http.Server

const TIMEOUT = 10 // timeout for server initialization
var connected = 0

func Initialize(landingPagePath string) error {
	// step: serve event server
	conn := event.CreateServerConnector()
	if err := conn.Open(); err != nil {
		return errors.SystemNew("Failed to initialize intaller server.", err)
	}
	global.Connector = conn
	server.InitSystemService(conn, nil)

	// step: init web server
	log.Println("Listening and serve @port", PORT, "www dir =", landingPagePath)
	webserver = &http.Server{Addr: ":" + PORT}
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		connected++
		http.FileServer(http.Dir(landingPagePath)).ServeHTTP(rw, r)
	})
	// step: listen and serve to target port
	go func() {
		elapsed := time.Now().Unix()
		for {
			err := webserver.ListenAndServe()
			if err == nil {
				break
			}
			// step: check if this is waiting for too long
			diff := time.Now().Unix() - elapsed
			if diff > TIMEOUT {
				log.Println("Server error.", err.Error())
				break
			} else {
				log.Println("Issue found, retrying...", err.Error())
			}
			// sleep for a while
			time.Sleep(time.Millisecond * 500)
		}
	}()
	conn.SetStatus(system.UPDATING, nil)
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
	log.Println("Shutting down event and server...")
	// close event server
	if global.Connector != nil {
		if err := global.Connector.Close(); err != nil {
			log.Println("Error closing connector.", err.Error())
		}
	}
	return webserver.Shutdown(context.TODO())
}

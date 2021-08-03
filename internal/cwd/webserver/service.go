package main

import (
	"context"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	_ "net/http/pprof"
)

const PORT = "80"
const PAGE_LANDING = constants.SYSTEM_SERVICE_ID
const PAGE_COMPANIONAPP = "ela.companion"

type MyService struct {
	server     *http.Server
	running    bool
	fileServer http.Handler
}

func (s *MyService) OnStart() error {
	s.running = true
	s.server = &http.Server{Addr: ":" + PORT}
	wwwPath := path.GetSystemWWW()
	fsrv, err := s.getFileserver(PAGE_LANDING)
	if err != nil {
		return errors.SystemNew("Unable to find the landing page", err)
	}
	s.fileServer = fsrv
	lastPkg := PAGE_LANDING

	// handle any requests
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		pkg := PAGE_LANDING
		log.Println("serve", url)
		// retrieve the package based from url
		if url != "/" {
			splits := strings.Split(url, "/")
			pkg = splits[1]
		}
		// switch package?
		if pkg != lastPkg {
			fileserver, err := s.getFileserver(pkg)
			if err == nil {
				log.Println("Package", pkg, "selected")
				lastPkg = pkg
				s.fileServer = fileserver
				debug.FreeOSMemory()
			}
		}

		// does the file exist? then serve the file
		if _, err := os.Stat(wwwPath + url); err == nil {
			s.fileServer.ServeHTTP(rw, r)
		} else { // hence use the index file
			http.ServeFile(rw, r, wwwPath+"/"+lastPkg)
		}
	})
	go func() {
		log.Println("Start listening to " + PORT)
		if err := s.server.ListenAndServe(); err != nil {
			log.Println("Server error", err.Error())
			s.running = false
		}
	}()
	return nil
}

func (s *MyService) getFileserver(packageId string) (http.Handler, error) {
	loc := path.GetSystemWWW() + "/" + packageId
	if _, err := os.Stat(loc); err == nil {
		return http.FileServer(http.Dir(loc)), nil
	} else {
		return nil, err
	}
}

func (s *MyService) OnEnd() error {
	return s.server.Shutdown(context.TODO())
}

func (s MyService) IsRunning() bool {
	return s.running
}

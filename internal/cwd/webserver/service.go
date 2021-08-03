package main

import (
	"context"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/internal/cwd/webserver/fs"
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
	dirHandler *fs.Dir
}

func (s *MyService) OnStart() error {
	s.running = true
	s.server = &http.Server{Addr: ":" + PORT}
	wwwPath := path.GetSystemWWW()
	s.dirHandler = &fs.Dir{}
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
		//log.Println("serve", url)
		// retrieve the package based from url
		if url != "/" {
			splits := strings.Split(url, "/")
			tmp := splits[1]
			// is this a package?
			if _, err := os.Stat(wwwPath + "/" + tmp); err == nil {
				pkg = tmp
				r.URL.Path = "/"
				if len(splits) > 1 {
					r.URL.Path = strings.Join(splits[2:], "/")
				}
			}
		}
		//log.Println(pkg, r.URL.Path)
		// switch package?
		if pkg != lastPkg {
			//fileserver, err := s.getFileserver(pkg)
			loc := path.GetSystemWWW() + "/" + pkg
			log.Println("Package", pkg, "selected")
			lastPkg = pkg
			s.dirHandler.SetPath(loc)
			debug.FreeOSMemory()
		}
		// does the file exist? then serve the file
		//if _, err := os.Stat(wwwPath + url); err == nil {
		s.fileServer.ServeHTTP(rw, r)
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
		s.dirHandler.CurrentPath = loc
		return http.FileServer(s.dirHandler), nil
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

package web

import (
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/path"

	_ "net/http/pprof"
)

const PORT = "80"
const PAGE_LANDING = constants.SYSTEM_SERVICE_ID
const PAGE_COMPANIONAPP = "ela.companion"

type WebService struct {
	running bool
}

func (s *WebService) Start() error {
	s.running = true
	wwwPath := path.GetSystemWWW()
	lastPkg := PAGE_LANDING

	// handle any requests
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		pkg := lastPkg
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
					r.URL.Path += strings.Join(splits[2:], "/")
				}
			}
		} else {
			pkg = PAGE_LANDING
		}
		//log.Println(pkg, r.URL.Path)
		// switch package?
		if pkg != lastPkg {
			log.Println("Package", pkg, "selected")
			lastPkg = pkg
			debug.FreeOSMemory()
		}
		fpath := wwwPath + "/" + lastPkg + r.URL.Path
		f, err := os.Stat(fpath)
		if err != nil || f.IsDir() {
			fpath = wwwPath + "/" + lastPkg + "/index.html"
		}
		//log.Println(fpath)
		http.ServeFile(rw, r, fpath)
	})

	return nil
}

func (s *WebService) Close() error {
	return nil
}

func (s WebService) IsRunning() bool {
	return s.running
}

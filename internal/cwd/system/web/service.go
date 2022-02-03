// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

// These file provides functionality for http client requests

package web

import (
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/event/data"
	"github.com/cansulting/elabox-system-tools/foundation/path"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/appman"
	"github.com/cansulting/elabox-system-tools/internal/cwd/system/global"
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
			s.onPackageSelected(pkg)
			lastPkg = pkg
			runtime.GC()
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

// callback when package was selected
func (s *WebService) onPackageSelected(pkg string) {
	if pkg == PAGE_LANDING {
		return
	}
	global.Logger.Debug().Str("category", "web").Msg("Package " + pkg + " selected.")
	// start the activity
	if err := appman.LaunchAppActivity(pkg, nil, data.NewActionById(constants.ACTION_APP_LAUNCH)); err != nil {
		global.Logger.Error().Err(err).Str("category", "web").Msg("Failed to launch activity.")
	}
}

func (s *WebService) Close() error {
	return nil
}

func (s WebService) IsRunning() bool {
	return s.running
}

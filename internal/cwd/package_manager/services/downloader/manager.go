package downloader

import (
	"errors"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
)

// data for resumable downloads
type resumableData struct {
	Id   string
	Src  string
	Dest string
}

var downloadTasks = make(map[string]*Task)
var isInit = false
var wwwImages = ""

//var resumableDl = make(map[string]resumableData, 0)

const RETRY_DOWNLOAD_COUNT = 3

const RETRY_DOWNLOAD_DELAY = 2 // in seconds

func initialize() {
	if isInit {
		return
	}
	isInit = true
	if err := os.MkdirAll(global.DownloadCache, perm.PUBLIC_WRITE); err != nil {
		logger.GetInstance().Error().Msg("unable to create cache directory " + global.DownloadCache)
	}
}

func AddDownload(id string, url string, mode DownloadMode) *Task {
	initialize()
	task := NewTask(id, url, global.DownloadCache+"/"+id, mode)
	downloadTasks[id] = task
	// task.onStateChanged = onTaskStateChanged
	// task.onError = onTaskError
	//task.Start()
	return task
}

func GetDownload(id string) *Task {
	return downloadTasks[id]
}

func RemoveDownload(id string) {
	delete(downloadTasks, id)
}

// this cache an image relative to www path
func AddCacheImage(id string, src string, dest string) error {
	initialize()
	addToResumables(id, src, dest)
	task := NewTask(id, src, dest, HTTP)
	retried := -1
	success := true
	if err := os.MkdirAll(path.Dir(dest), perm.PUBLIC_WRITE); err != nil {
		logger.GetInstance().Error().Err(err).Caller().Msg("failed to create parent dir for cache " + id)
	}
	// retry for failure download
	for retried < RETRY_DOWNLOAD_COUNT {
		retried++
		success = true
		if err := task.Start(); err != nil {
			success = false
			logger.GetInstance().Error().Err(err).Caller().Msg("failed to download cache " + src)
			if retried < RETRY_DOWNLOAD_COUNT {
				time.Sleep(time.Second * RETRY_DOWNLOAD_DELAY)
			}
		}
	}
	if success {
		removeFromResumables(id)
	} else {
		return errors.New("failed to download ")
	}
	return nil
}

func IdentifyDownloadMode(downloadUrl string) DownloadMode {
	if strings.HasPrefix(downloadUrl, "http://") ||
		strings.HasPrefix(downloadUrl, "https://") {
		return HTTP
	}
	return IPFS
}

func addToResumables(id string, src string, dest string) {
	// resumableDl[id] = resumableData{
	// 	Id:   id,
	// 	Src:  src,
	// 	Dest: dest,
	// }
}

func removeFromResumables(id string) {
	// delete(resumableDl, id)
}

// func onTaskStateChanged(task *Task) {
// 	// do something
// }

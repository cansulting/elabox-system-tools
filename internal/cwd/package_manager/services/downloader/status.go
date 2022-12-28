package downloader

type Status int

// the current status of download. 0 mean idle, 1 mean downloading, 2 mean paused, 3 mean stopped
const (
	Idle        Status = 0
	Downloading Status = 1
	Paused      Status = 2
	Stopped     Status = 3
	Error       Status = 4
	Finished    Status = 5
)

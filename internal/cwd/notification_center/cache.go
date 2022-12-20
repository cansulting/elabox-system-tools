package main

import "github.com/cansulting/elabox-system-tools/foundation/path"

var notifQueue []NotifData

const FILE_NAME = "notif.cache"

// cache file location
func getCachePath() string {
	return path.GetSystemAppDirData(PACKAGE_ID) + "/" + FILE_NAME
}

func initQueue() error {
	// already loaded? then skip
	if notifQueue != nil {
		return nil
	}
	notifQueue = make([]NotifData, 0, NOTIF_QUEUE_LIMIT)
	// sample
	notifQueue = append(notifQueue, NotifData{Title: "Test", Message: "This is sample", Status: Unread})
	// load notification here
	return nil
}

func saveQueue() error {

	return nil
}

func AddNotif(data NotifData) error {
	if err := initQueue(); err != nil {
		return err
	}
	// is reach queue limit? then dequeue the old data
	if len(notifQueue) >= NOTIF_QUEUE_LIMIT {
		notifQueue = notifQueue[1:]
	}
	notifQueue = append(notifQueue, data)
	return nil
}

// retrieve queue of
// @page: start with 1
func RetrieveNotif(page uint, length uint) ([]NotifData, error) {
	if err := initQueue(); err != nil {
		return nil, err
	}
	// if not enough then just return all data
	totalNotif := len(notifQueue)
	if totalNotif-int(length) <= 0 {
		return notifQueue, nil
	}

	// not within the limit? return empty
	startI := (totalNotif - (int(page) * int(length)))
	endI := startI + int(length)
	if startI < 0 {
		return nil, nil
	}

	result := notifQueue[startI:endI]
	return result, nil
}

package main

import (
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

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
	//notifQueue = append(notifQueue, NotifData{Title: "Test", Message: "This is sample", Status: Unread})
	// load notification here
	// content, err := ioutil.ReadFile(FILE_NAME)
	// if err != nil {
	// 	//log.Fatal(err)
	// 	fmt.Println("File Does Not Exist")
	// 	return err
	// }
	// err = json.Unmarshal(content, &notifQueue)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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

	// // add notif
	// notifQueue, err := json.Marshal(notifQueue)
	// if err != nil {
	// 	panic(err)
	// }
	// err = ioutil.WriteFile(FILE_NAME, notifQueue, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
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
	if page <= 0 {
		return nil, nil
	}
	// not within the limit? return empty
	endI := (totalNotif - (int(page-1) * int(length)))
	if endI <= 0 {
		return nil, nil
	}
	startI := endI - int(length)
	if startI < 0 {
		startI = 0
	}
	result := notifQueue[startI:endI]
	return result, nil
}

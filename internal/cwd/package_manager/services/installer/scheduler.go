// this file manages the schedule for each installation
// Packages can download simultaneously, but the installer can only install one package at a time.
// If a package B downloaded successfully, it needs to wait the current package A being installed.
// Then the installer can install package B.

package installer

import "github.com/cansulting/elabox-system-tools/foundation/logger"

var currentSchedule *Task // current scheduled task
var installQueue []*Task = make([]*Task, 0, 5)

// use to add task to install schedule
func addToSchedule(task *Task) {
	logger.GetInstance().Debug().Msg(task.Id + " added to install schedule")
	installQueue = append(installQueue, task)
	if !hasCurrentInstalling() {
		startQueue()
	}
}

// check if theres currently installing
func hasCurrentInstalling() bool {
	return currentSchedule != nil
}

// use to start the next queued task
func startQueue() {
	if len(installQueue) == 0 {
		return
	}
	currentSchedule = installQueue[0]
	installQueue = installQueue[1:]
	currentSchedule.install("")
}

// finish the current scheduled task and start the next one
func finishCurrentSchedule() {
	if currentSchedule == nil {
		return
	}
	logger.GetInstance().Debug().Msg(currentSchedule.Id + " install finished")
	currentSchedule = nil
	startQueue()
}

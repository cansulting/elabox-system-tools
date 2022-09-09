package main

import "github.com/cansulting/elabox-system-tools/internal/cwd/events/listeners"

func main() {
	listeners.ListenToShutdown()
}
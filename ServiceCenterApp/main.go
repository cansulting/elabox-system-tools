package main

import (
	"time"
)

var running bool = true

func main() {
	// this runs the server
	RunServer()
	for running {
		time.Sleep(time.Second * 1)
	}
	// closes the connection
	Close()
}

package main

import (
	"testing"
	"time"
)

func TestService(test *testing.T) {
	s := MyService{}
	s.OnStart()
	for s.IsRunning() {
		go time.Sleep(time.Second)
	}
}

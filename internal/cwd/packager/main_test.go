package main

import "testing"

func TestPackage(t *testing.T) {
	if err := load("../../builds/linux/packager/companion.json"); err != nil {
		t.Error(err)
	}
}

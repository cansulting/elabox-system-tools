package main

import "testing"

func TestPackage(t *testing.T) {
	if err := load("../../builds/windows/packager/packageinstaller.json"); err != nil {
		t.Error(err)
	}
}

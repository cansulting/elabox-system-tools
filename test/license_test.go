package main

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

const host = "208.87.134.80:1236" //staging

func retrieveSerial() string {
	fname := "/proc/cpuinfo"
	bytes, _ := os.ReadFile(fname)
	splits := strings.Split(string(bytes), "\n")
	for _, split := range splits {
		keypair := strings.Split(split, ":")
		if len(keypair) > 0 {
			key := strings.Trim(keypair[0], "\t")
			if key == "Serial" {
				return strings.Trim(keypair[1], " ")
			}
		}
	}
	return ""
}

// use to check if device was already registered
func TestRegistrationChecker(test *testing.T) {
	serial := retrieveSerial()
	res, err := http.Post("http://"+host+"/apiv1/rewards/check-device?serial="+serial, "application/json", nil)
	if err != nil {
		test.Error(res)
		return
	}
	contents, err := io.ReadAll(res.Body)
	if err != nil {
		test.Error(err)
		return
	}
	test.Log(string(contents))
}

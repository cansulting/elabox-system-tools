package main

import (
	"testing"

	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/ipfs"
)

const TEST_CID = "QmePfgfoB27qQyWEV2oJNQMQkeXit1dCEue3WJHU85fHUE"

func Test_DownloadJson(t *testing.T) {
	var output map[string]interface{}
	err := ipfs.DownloadJson(TEST_CID, &output)
	if err != nil {
		t.Error(err)
	}
	println(output)
}

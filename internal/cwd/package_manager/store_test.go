package main

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/data"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/global"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/downloader"
	"github.com/cansulting/elabox-system-tools/internal/cwd/package_manager/services/installer"
)

const TEST_PKG_PATH = "/../build/ela.sample/ela.sample.box"

var INSTALL_DEF = data.InstallDef{Id: "ela.sample", Url: "https://storage.googleapis.com/elabox-staging/packages/6.box"}

func Test_RetrieveListing(t *testing.T) {
	logger.Init("ela.store.test")
	// step: retrieve store listing
	// if err := store_lister.CheckUpdates(); err != nil {
	// 	t.Error("unable to retrieve store listing. inner: " + err.Error())
	// }
}

// test for retrieving apps information and states
func Test_RetrieveAppsState(t *testing.T) {
	logger.Init("ela.store.test")
	pkgs, err := RetrieveAllApps(true)
	if err != nil {
		t.Error("unable to retrieve all installed packages. inner: " + err.Error())
		return
	}
	t.Log("retrieved all installed packages " + strconv.Itoa(len(pkgs)))
}

// test for retrieving specific app detailed information
func Test_RetrieveAppDetail(t *testing.T) {
	logger.Init("ela.store.test")
	pkgs, err := RetrieveAllApps(true)
	if err != nil {
		t.Error("unable to retrieve all installed packages. inner: " + err.Error())
		return
	}
	pkg := pkgs[0]
	pkgDetail, err := RetrieveApp(pkg.Id, "")
	if err != nil {
		t.Error("unable to retrieve app detail. inner: " + err.Error())
		return
	}
	t.Log("retrieved app detail: " + pkgDetail.Name)
}

// install app test
func Test_InstallPackage(t *testing.T) {
	logger.Init("ela.store.test")
	task, err := installer.CreateInstallTask(data.InstallDef{}, nil)
	if err != nil {
		t.Error("unable to create install task. inner: " + err.Error())
		return
	}
	handler, err := rpc.NewRPCHandlerDefault()
	global.RPC = handler
	if err != nil {
		t.Error("unable to create rpc handler. inner: " + err.Error())
		return
	}
	path, _ := os.Getwd()
	path += TEST_PKG_PATH
	if err := task.StartFromFile(path); err != nil {
		t.Error("unable to install package. inner: " + err.Error())
		return
	}
	time.Sleep(time.Second * 3)
}

// test for downloading app
func Test_DownloadPackageHttp(t *testing.T) {
	logger.Init("ela.store.test")
	url := "https://storage.googleapis.com/elabox-staging/packages/6.box"
	savePath := "./sample.box"
	task := downloader.NewTask("sample", url, savePath, downloader.HTTP)
	if err := task.Start(); err != nil {
		t.Error("unable to download package. inner: " + err.Error())
		return
	}
}

func Test_DownloadPackageIPFS(t *testing.T) {
	logger.Init("ela.store.test")
	url := TEST_CID
	savePath := "./sample.box"
	task := downloader.NewTask("web3._store", url, savePath, downloader.IPFS)
	if err := task.Start(); err != nil {
		t.Error("unable to download package. inner: " + err.Error())
		return
	}
}

// use to test installation for dependencies
// func Test_InstallWithDependencies(t *testing.T) {
// 	logger.Init("ela.store.test")
// 	handler, err := rpc.NewRPCHandlerDefault()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	global.RPC = handler
// 	if err := broadcast.Init(); err != nil {
// 		t.Error("failed to init broadcast", err)
// 		return
// 	}
// 	dependencies := []string{"trinity.pasar", "ipfs"}
// 	task := installer.CreateTask("ela.mainchain", TEST_CID, dependencies)
// 	task.Start()
// 	for {
// 		if task.Status == global.Installed {
// 			break
// 		}
// 		if task.ErrorCode != 0 {
// 			t.Error("Failed installing with dependencies with error code", task.ErrorCode)
// 			break
// 		}
// 	}
// }

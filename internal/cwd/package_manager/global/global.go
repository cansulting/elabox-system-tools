package global

import (
	"github.com/cansulting/elabox-system-tools/foundation/app"
	"github.com/cansulting/elabox-system-tools/foundation/app/rpc"
	"github.com/cansulting/elabox-system-tools/foundation/path"
)

const PackageId = "ela.pkgmgr"
const InstallerId = "ela.installer"
const ACCOUNT_PACKAGE_ID = "ela.account"

// FLAGS and Constants
const RetrieveSystem_Delay = 120

// path and caches
const CacheDir = path.PATH_SYSTEM_DATA + "/" + PackageId + "/cache"
const StoreCache = CacheDir + "/packages.cache"
const DownloadCache = CacheDir + "/downloads"

var AppController *app.Controller
var RPC *rpc.RPCHandler

// actions
const RETRIEVE_PACKAGES = PackageId + ".actions.RETRIEVE_PACKAGES"
const RETRIEVE_PACKAGE = PackageId + ".actions.RETRIEVE_PACKAGE"
const INSTALL_PACKAGE = PackageId + ".actions.INSTALL_PACKAGE"
const UNINSTALL_PACKAGE = PackageId + ".actions.UNINSTALL_PACKAGE"
const CANCEL_INSTALL_PACKAGE = PackageId + ".actions.CANCEL_INSTALL_PACKAGE"
const AC_DOWNLOAD_UPDATE = PackageId + ".actions.DOWNLOAD_UPDATE"
const AC_INSTALL_UPDATE = PackageId + ".actions.INSTALL_UPDATE"

// broadcasts
const UPDATE_AVAILABLE = PackageId + ".broadcast.UPDATE_AVAILABLE"
const INSTALL_PROGRESS = PackageId + ".broadcast.INSTALL_PROGRESS"
const BROADCAST_ERROR = PackageId + ".broadcast.ERROR"
const INSTALL_STATE = PackageId + ".broadcast.INSTALL_STATE"
const RETRIEVE_SYS_VERSION = PackageId + ".actions.RETRIEVE_SYSTEM_VERSION"

// installer actions
const INSTALLER_PROGRESS = "ela.installer.broadcast.PROGRESS"
const INSTALLER_STATE_CHANGE = "ela.installer.broadcast.STATE_CHANGED"
const INSTALLER_ERROR = "ela.installer.broadcast.ERROR"

const SYS_INFO_PATH = "/usr/ela/system/ela.system/info.json"

// CACHES
const IMAGES_CACHE_PATH = "/_images" // cache path for all package images
const STORE_CACHE_PATH = CacheDir + "/stores.cache"
const STOREHUB_CACHE_PATH = CacheDir + "/storehub.cache"

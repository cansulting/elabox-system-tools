#!/bin/bash
# unset go lang env variables
go env -u GOOS
go env -u GOARCH
go env -u GO111MODULE
go env -u CC
go env -u CXX

PROJ_HOME=../../..
ELA_SRC=$PROJ_HOME/Elastos.ELA
EID_SRC=$PROJ_HOME/Elastos.ELA.SideChain.EID
ESC_SRC=$PROJ_HOME/Elastos.ELA.SideChain.ESC
GLIDE_SRC=$PROJ_HOME/glide-frontend
ELA_LANDING=$PROJ_HOME/elabox-companion-landing
ELA_REWARDS=$PROJ_HOME/elabox-rewards
ELA_LOGS=$PROJ_HOME/elabox-logs
ELA_STORE=$PROJ_HOME/elabox-dapp-store
ELA_SETUP=$PROJ_HOME/elabox-setup-wizard
ELA_DASHBOARD=$PROJ_HOME/elabox-dashboard
cos=$(go env GOOS)                  # current os. 
carc=$(go env GOARCH)               # current archi
packageinstaller=packageinstaller           # package installer project name
system_name=system                  # system project name
packager=packager             
target=$cos
arch=$carc
gobuild='go build -tags DEBUG'    # build command
MODE=DEBUG

# FLAGS
while getopts o:a:d flag
do
    case "${flag}" in
        o) target=${OPTARG};;
        a) arch=${OPTARG};;
    esac
done
echo "Optional commandline params -o(target) -a(arch))"
echo "eg. -o linux -a arm64"
echo "OS="$target
echo "Arch="$arch

# release mode?
echo "Build Mode: 1 - RELEASE, 2 - STAGING, 3 - Default - DEBUG (leave empty if DEBUG)"
read mode
if [ "$mode" == "1" ]; then
    MODE=RELEASE
    gobuild='go build -ldflags "-w -s" -tags RELEASE'
elif [ "$mode" == "2" ]; then
    MODE=STAGING
    gobuild='go build -tags STAGING'
fi
echo "Mode=$MODE"
echo ""

# where binaries will be saved
go env -w CGO_ENABLED=1
echo "cgo enabled"

#####################
# Questions
#####################
echo "Rebuild elabox landing page? (y/n)"
read answerLand
echo "Rebuild elastos binaries? (y/n)"
read answerEla
if [ -d "$ELA_LOGS" ]; then 
    echo "Rebuild logging service? (y/n)"
    read answerLog
fi
echo "Rebuild Glide? (y/n)"
read answerGlide
echo "Rebuild elastos dapp store? (y/n)"
read answerDstore
echo "Rebuild Setup Wizard? (y/n)"
read answerSetup
echo "Rebuild Dashboard? (y/n)"
read answerDashboard

#####################
# build packager
#####################
buildpath=../builds/$target
echo "Building " $packager
mkdir -p $buildpath/packager
eval "$gobuild" -o $buildpath/packager/$packager ../cwd/$packager

ln -sf $PWD/$buildpath/packager/$packager /usr/local/bin/$packager

#####################
# build system binaries
#####################
if [ "$target" == "linux" ]; then
    if [ "$arch" == "arm64" ]; then
        # specific gcc for arm64
        go env -w CC=aarch64-linux-gnu-gcc
    elif [ "$arch" == "386" ]; then
        # intel
        go env -w CC=i686-linux-gnu-gcc
    fi
elif [ "$target" == "windows" ]; then 
    # windows intel
    if [ "$arch" == "386" ]; then
        go env -w CXX=i686-w64-mingw32-g++ 
        go env -w CC=i686-w64-mingw32-gcc
    # windows amd
    else
        go env -w CXX=x86_64-w64-mingw32-g++ 
        go env -w CC=x86_64-w64-mingw32-gcc
    fi
fi

go env -w GOOS=$target 
go env -w GOARCH=$arch
echo "Building " $packageinstaller
mkdir -p $buildpath/$packageinstaller/bin
eval "$gobuild" -o $buildpath/$packageinstaller/bin ../cwd/$packageinstaller
ln -sf $PWD/$buildpath/$packageinstaller/bin/$packageinstaller /usr/local/bin/$packageinstaller
echo "Building Elabox System"
eval "$gobuild" -o $buildpath/$system_name/bin ../cwd/$system_name
programName=$(jq ".program" $buildpath/$system_name/info.json | sed 's/\"//g')
mv $buildpath/$system_name/bin/$system_name $buildpath/$system_name/bin/$programName 

# build account manager
echo "Building Account Manager"
mkdir -p $buildpath/account_manager/bin
eval "$gobuild" -o $buildpath/account_manager/bin ../cwd/account_manager
programName=$(jq ".program" $buildpath/account_manager/info.json | sed 's/\"//g')
mv $buildpath/account_manager/bin/account_manager $buildpath/account_manager/bin/$programName 

# build notification
echo "Building Notification System"
mkdir -p $buildpath/notification_center/bin
eval "$gobuild" -o $buildpath/notification_center/bin ../cwd/notification_center
programName=$(jq ".program" $buildpath/notification_center/info.json | sed 's/\"//g')
mv $buildpath/notification_center/bin/notification_center $buildpath/notification_center/bin/$programName 

# build package manager
echo "Building Package Manager"
mkdir -p $buildpath/package_manager/bin
eval "$gobuild" -o $buildpath/package_manager/bin ../cwd/package_manager
programName=$(jq ".program" $buildpath/package_manager/info.json | sed 's/\"//g')
mv $buildpath/package_manager/bin/package_manager $buildpath/package_manager/bin/$programName 

# build reward if exists
if [ -d "$ELA_REWARDS" ]; then 
    wd=$PWD
    cd $ELA_REWARDS/scripts
    ./build.sh -o $target -a $arch -d $MODE
    cd $wd
fi

# build app logs
if [ "$answerLog" == "y" ]; then 
    wd=$PWD
    cd $ELA_LOGS/scripts
    ./build.sh -o $target -a $arch -d $MODE
    cd $wd
fi

# unset env variables
go env -u CC
go env -u CXX

##################################
# build system landing page
##################################
if [ "$answerLand" == "y" ]; then
    wd=$PWD
    cd $ELA_LANDING
    sudo npm install
    sudo npm run build
    cd $wd
    rm -r $buildpath/system/www 
    mkdir -p $buildpath/system/www
    cp -r $ELA_LANDING/build/* $buildpath/system/www
fi

##################################
# elastos mainchain, eid, cli, esc
##################################
buildELA() {
    echo "Building ELA from source..."
    wd=$PWD
    cd $ELA_SRC 
    make all
    echo "Done building Ela."
}
if [ "$answerEla" == "y" ]; then
    buildELA
    echo "Building EID from source..."
    go env -w GO111MODULE=off
    cd $wd
    cd $EID_SRC
    make geth
    cd $wd
    echo "Building ESC from source..."
    cd $wd
    cd $ESC_SRC
    make geth
    cd $wd
    go env -u GO111MODULE

    targetdir=$buildpath
    echo "Copying mainchain, eid and cli "
    # mainchain
    mainchainlib=$buildpath/mainchain/bin
    mkdir -p $mainchainlib
    cp $ELA_SRC/ela-cli $mainchainlib
    cp $ELA_SRC/ela $mainchainlib
    chmod +x $mainchainlib/ela $mainchainlib/ela-cli
    # eid
    eidbin=$buildpath/eid/bin
    mkdir -p $eidbin
    cp ${EID_SRC}/build/bin/geth $eidbin
    chmod +x $eidbin/geth
    mv $eidbin/geth $eidbin/ela.eid
    # esc
    escbin=$buildpath/esc/bin
    mkdir -p $escbin
    cp ${ESC_SRC}/build/bin/geth $escbin
    chmod +x $escbin/geth
    mv $escbin/geth $escbin/ela.esc
    # carrier
    carrierlib=$buildpath/carrier/bin
    chmod +x $carrierlib/ela-bootstrapd

    packager $buildpath/eid/packager.json
    packager $buildpath/esc/packager.json
    packager $buildpath/carrier/packager.json
    packager $buildpath/feeds/packager.json
    packager $buildpath/mainchain/packager.json
fi

#########################
# build Glide?
#########################
if [ "$answerGlide" == "y" ]; then
    echo "Building Glide..."
    wd=$PWD
    cd $GLIDE_SRC
   # sudo npm install
    #sudo npm run build
    cd $wd
    rm -r $buildpath/glide/www && mkdir -p $buildpath/glide/www
    cp -r $GLIDE_SRC/build/* $buildpath/glide/www
    packager $buildpath/glide/packager.json
fi

#########################
# build dapp store?
#########################
if [ "$answerDstore" == "y" ]; then
    wd=$PWD
    cd $ELA_STORE/scripts
    ./build.sh -o $target -a $arch -d $MODE
    cd $wd
fi

#########################
# build setup wizard?
#########################
if [ "$answerSetup" == "y" ]; then
    wd=$PWD
    cd $ELA_SETUP/scripts
    ./build.sh -o $target -a $arch -d $MODE
    cd $wd
fi

#########################
# build dashboard?
#########################
if [ "$answerDashboard" == "y" ]; then
    wd=$PWD
    cd $ELA_DASHBOARD/scripts
    ./build.sh -o $target -a $arch -d $MODE
    cd $wd
fi

#########################
# Packaging
#########################
echo "Start packaging..."
packager $buildpath/$packageinstaller/packager.json
packager $buildpath/account_manager/packager.json
packager $buildpath/notification_center/packager.json
packager $buildpath/package_manager/packager.json
packager $buildpath/$system_name/packager.json

go env -u CGO_ENABLED

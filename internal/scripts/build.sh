#!/bin/bash
PROJ_HOME=../../..
ELA_NODES=$PROJ_HOME/elabox-binaries/binaries
ELA_SRC=$PROJ_HOME/Elastos.ELA
EID_SRC=$PROJ_HOME/Elastos.ELA.SideChain.EID
ESC_SRC=$PROJ_HOME/Elastos.ELA.SideChain.ESC
ELA_COMPANION=$PROJ_HOME/elabox-companion
ELA_LANDING=$PROJ_HOME/landing-page
ELA_REWARDS=$PROJ_HOME/elabox-rewards
ELA_LOGS=$PROJ_HOME/elabox-logs
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
echo "Rebuild companion client & server? 1 - All, 2 - Client, 3 - Server, Enter - none"
read answerComp
echo "Rebuild elabox landing page? (y/n)"
read answerLand
echo "Do you want to rebuild elastos binaries? (y/n)"
read answerEla

#####################
# build packager
#####################
go env -u GOOS
go env -u GOARCH
buildpath=../builds/$target
echo "Building " $packager
mkdir -p $buildpath/packager
eval "$gobuild" -o $buildpath/packager/$packager ../cwd/$packager
ln -sf $PWD/$buildpath/packager/$packager /bin/$packager

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
else
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
ln -sf $PWD/$buildpath/$packageinstaller/bin/$packageinstaller /bin/$packageinstaller
echo "Building Elabox System"
eval "$gobuild" -o $buildpath/$system_name/bin ../cwd/$system_name
programName=$(jq ".program" $buildpath/$system_name/info.json | sed 's/\"//g')
mv $buildpath/$system_name/bin/$system_name $buildpath/$system_name/bin/$programName 

# build reward if exists
if [ -d "$ELA_REWARDS" ]; then 
    wd=$PWD
    cd $ELA_REWARDS/scripts
    ./build.sh -o $target -a $arch -d $MODE
    cd $wd
fi

# build app logs
wd=$PWD
cd $ELA_LOGS/scripts
./build.sh -o $target -a $arch -d $MODE
cd $wd

# unset env variables
go env -u CC
go env -u CXX


#########################
# build companion app?
#########################
built=0
# client building
if [[ "$answerComp" == "1" || "$answerComp" == "2" ]]; then 
    echo "Start building client companion app, please wait this will take awhile..." 
    initDir=$PWD
    cd $ELA_COMPANION/src_client
    sudo npm install
    sudo npm run build
    cd $initDir
    rm -r $buildpath/companion/www
    mkdir -p $buildpath/companion/www
    cp -r $ELA_COMPANION/src_client/build/* $buildpath/companion/www
    built=1
fi
# server building
if [[ "$answerComp" == "1" || "$answerComp" == "3" ]]; then 
    echo "Start building server companion app, please wait this will take awhile..." 
    initDir=$PWD
    cd $ELA_COMPANION/src_server
    sudo npm install
    cd $initDir
    mkdir -p $buildpath/companion/nodejs
    cp -r $ELA_COMPANION/src_server/* $buildpath/companion/nodejs
    built=1
fi
if [ "$built" == "1" ]; then 
    echo "Build success!"
    echo "Packaging..."
    packager $buildpath/companion/packager.json
fi

##################################
# build landing page
##################################
if [ "$answerLand" == "y" ]; then
    wd=$PWD
    cd $ELA_LANDING
    sudo npm install
    sudo npm run build
    cd $wd
    cp -r $ELA_LANDING/build/* $buildpath/system/www
fi

##################################
# elastos mainchain, eid, cli, esc
##################################
if [ "$answerEla" == "y" ]; then
    echo "Building ELA from source..."
    wd=$PWD
    cd $ELA_SRC 
    make all
    echo "Building EID from source..."
    cd $wd
    cd $EID_SRC
    make geth
    cd $wd
    echo "Building ESC from source..."
    cd $wd
    cd $ESC_SRC
    make geth
    cd $wd

    targetdir=$buildpath
    echo "Copying mainchain, eid and cli @$ELA_NODES"
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
    # esc
    escbin=$buildpath/esc/bin
    mkdir -p $escbin
    cp ${ESC_SRC}/build/bin/geth $escbin
    chmod +x $escbin/geth
    mv $escbin/geth $escbin/esc
    # carrier
    carrierlib=$buildpath/carrier/bin
    mkdir -p $carrierlib
    cp ${ELA_NODES}/ela-bootstrapd $carrierlib
    cp ${ELA_NODES}/bootstrapd.conf $carrierlib
    chmod +x $carrierlib/ela-bootstrapd
fi

#########################
# Packaging
#########################
echo "Start packaging..."
packager $buildpath/eid/packager.json
packager $buildpath/esc/packager.json
packager $buildpath/carrier/packager.json
packager $buildpath/mainchain/packager.json
packager $buildpath/$packageinstaller/packager.json
packager $buildpath/feeds/packager.json
packager $buildpath/$system_name/packager.json

go env -u CGO_ENABLED

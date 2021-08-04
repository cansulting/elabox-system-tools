#!/bin/bash
PROJ_HOME=~
ELA_NODES=$PROJ_HOME/elabox-binaries/binaries
ELA_COMPANION=$PROJ_HOME/elabox-companion
ELA_LANDING=$PROJ_HOME/landing-page
cos=$(go env GOOS)                  # current os. 
carc=$(go env GOARCH)               # current archi
pkg_name=packageinstaller           # package installer project name
system_name=system                  # system project name
packager=packager
webserver=webserver                   
target=$cos
arch=$carc
gobuild='go build'                  # build command

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
echo "Release mode (y/n)?"
read mode
if [ "$mode" == "y" ]; then
    echo "Release=Enabled"
    gobuild='go build -ldflags "-w -s" -tags RELEASE'
fi
echo ""

# where binaries will be saved
go env -w CGO_ENABLED=1
echo "cgo enabled"

#####################
# build packager
#####################
go env -u GOOS
go env -u GOARCH
buildpath=../builds/$target
echo "Building " $packager
mkdir -p $buildpath/packager
eval "$gobuild" -o $buildpath/packager/$packager ../cwd/$packager

#####################
# build binaries
#####################
if [ "$target" == "linux" ]; then
    if [ "$arch" == "arm64" ]; then
        # specific gcc for arm64
        go env -w CC=aarch64-linux-gnu-gcc
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
echo "Building " $pkg_name
mkdir -p $buildpath/$pkg_name/bin
eval "$gobuild" -o $buildpath/$pkg_name/bin ../cwd/$pkg_name
echo "Building " $webserver
eval "$gobuild" -o $buildpath/$webserver/bin ../cwd/$webserver
echo "Building " $system_name
eval "$gobuild" -o $buildpath/$system_name/bin ../cwd/$system_name
mv $buildpath/$system_name/bin/$system_name $buildpath/$system_name/bin/main 
# unset env variables
go env -u CC
go env -u CXX

#########################
# build companion app?
#########################
built=0
echo "Rebuild companion client & server? 1 - All, 2 - Client, 3 - Server, Enter - none"
read answer
# client building
if [[ "$answer" == "1" || "$answer" == "2" ]]; then 
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
if [[ "$answer" == "1" || "$answer" == "3" ]]; then 
    echo "Start building server companion app, please wait this will take awhile..." 
    initDir=$PWD
    cd $ELA_COMPANION/src_server
    sudo npm install
    sudo npm run build
    cd $initDir
    mkdir -p $buildpath/companion/nodejs
    cp -r $ELA_COMPANION/src_server/* $buildpath/companion/nodejs
    built=1
fi
if [ "$built" == "1" ]; then 
    echo "Build success!"
    echo "Packaging..."
    pkgerPath=../builds/$cos/$packager/$packager
    $pkgerPath $buildpath/companion/packager.json
fi

##################################
# build landing page
##################################
echo "Rebuild elabox landing page? (y/n)"
read answer
if [ "$answer" == "y" ]; then
    wd=$PWD
    cd $ELA_LANDING
    sudo npm install
    sudo npm run build
    cd $wd
    cp -r $ELA_LANDING/build/* $buildpath/system/www
fi

##################################
# elastos mainchain, did, cli
##################################
targetdir=$buildpath
echo "Copying mainchain, did and cli @$ELA_NODES"
# mainchain
mainchainlib=$buildpath/mainchain/bin
mkdir -p $mainchainlib
cp ${ELA_NODES}/ela $mainchainlib
cp ${ELA_NODES}/ela-cli $mainchainlib
chmod +x $mainchainlib/ela $mainchainlib/ela-cli
cp ${ELA_NODES}/ela_config.json $mainchainlib
mv $mainchainlib/ela_config.json $mainchainlib/config.json
# did
didlib=$buildpath/did/bin
mkdir -p $didlib
cp ${ELA_NODES}/did $didlib
chmod +x $didlib/did
cp ${ELA_NODES}/did_config.json $didlib
mv $didlib/did_config.json $didlib/config.json
# carrier
carrierlib=$buildpath/carrier/bin
mkdir -p $carrierlib
cp ${ELA_NODES}/ela-bootstrapd $carrierlib
cp ${ELA_NODES}/bootstrapd.conf $carrierlib
chmod +x $carrierlib/ela-bootstrapd

#########################
# Packaging
#########################
echo "Start packaging..."
pkgerPath=../builds/$cos/$packager/$packager
$pkgerPath $buildpath/did/packager.json
$pkgerPath $buildpath/carrier/packager.json
$pkgerPath $buildpath/mainchain/packager.json
$pkgerPath $buildpath/$pkg_name/packager.json
$pkgerPath $buildpath/$webserver/packager.json
$pkgerPath $buildpath/$system_name/packager.json

go env -u CGO_ENABLED

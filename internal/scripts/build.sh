#!/bin/bash
ELA_NODES=../../../elabox-binaries/binaries
ELA_COMPANION=../../../elabox-companion
cos=$(go env GOOS)                  # current os. 
carc=$(go env GOARCH)               # current archi
pkg_name=packageinstaller           # package installer project name
system_name=system                  # system project name
packager=packager                   
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
output=$PWD/../builds/$target/bins
go env -w CGO_ENABLED=1
echo "cgo enabled"

#####################
# build packager
#####################
go env -u GOOS
go env -u GOARCH
echo "Building " $packager
eval "$gobuild" -o $output ../cwd/$packager

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
eval "$gobuild" -o $output ../cwd/$pkg_name
echo "Building " $system_name
eval "$gobuild" -o $output ../cwd/$system_name
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
    rm -r ../builds/$target/www/companion
    mkdir -p ../builds/$target/www/companion
    cp -r $ELA_COMPANION/src_client/build/* ../builds/$target/www/companion
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
    mkdir -p ../builds/$target/nodejs/companion
    cp -r $ELA_COMPANION/src_server/* ../builds/$target/nodejs/companion
    built=1
fi
if [ "$built" == "1" ]; then 
    echo "Build success!"
    echo "Packaging..."
    pkgerPath=../builds/$cos/bins/$packager
    $pkgerPath ../builds/$target/packager/companion.json
fi

##################################
# elastos mainchain, did, cli
##################################
targetdir=../builds/$target/libs
echo "Copying mainchain, did and cli @$ELA_NODES"
mkdir -p $targetdir/mainchain
mkdir -p $targetdir/did
mkdir -p $targetdir/carrier
# mainchain
cp ${ELA_NODES}/ela $targetdir/mainchain
cp ${ELA_NODES}/ela-cli $targetdir/mainchain
chmod +x $targetdir/mainchain/ela $targetdir/mainchain/ela-cli
cp ${ELA_NODES}/ela_config.json $targetdir/mainchain
mv $targetdir/mainchain/ela_config.json $targetdir/mainchain/config.json
# did
cp ${ELA_NODES}/did $targetdir/did
chmod +x $targetdir/did/did
cp ${ELA_NODES}/did_config.json $targetdir/did
mv $targetdir/did/did_config.json $targetdir/did/config.json
# carrier
cp ${ELA_NODES}/ela-bootstrapd $targetdir/carrier
cp ${ELA_NODES}/bootstrapd.conf $targetdir/carrier
chmod +x $targetdir/carrier/ela-bootstrapd
chmod 777 $targetdir/carrier/bootstrapd.conf

#########################
# Packaging
#########################
echo "Start packaging..."
pkgerPath=../builds/$cos/bins/$packager
$pkgerPath ../builds/$target/packager/did.json
$pkgerPath ../builds/$target/packager/carrier.json
$pkgerPath ../builds/$target/packager/mainchain.json
$pkgerPath ../builds/$target/packager/$pkg_name.json
$pkgerPath ../builds/$target/packager/$system_name.json

go env -u CGO_ENABLED

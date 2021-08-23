#!/bin/bash
target=$(go env GOOS) 
arch=$(go env GOARCH) 
pk=
gobuild='go build -tags DEBUG'    # build command
MODE=DEBUG

# FLAGS
while getopts o:a:p: flag
do
    case "${flag}" in
        o) target=${OPTARG};;
        a) arch=${OPTARG};;
        p) pk=${OPTARG};;
    esac
done
echo "Optional commandline params -o(target) -a(arch)) -p package"
echo "eg. -o linux -a arm64"
echo "OS="$target
echo "Arch="$arch

# release mode?
echo "Build Mode: 1 - RELEASE, 2 - STAGING, Default - DEBUG"
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

# build package
buildpath=../builds/$target/$pk
pkinfo=$buildpath/info.json
pkid=$(jq '.packageId' $pkinfo | sed 's/"//g')
echo "Package ID" $pkid
eval "$gobuild" -o $buildpath/bin/main ../cwd/$pk

# box the binary
../builds/$target/packager/packager $buildpath/packager.json

# install package
sudo ../builds/$target/packageinstaller/bin/packageinstaller ${buildpath}/${pkid}.box
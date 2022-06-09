#!/bin/bash
# upload the system package to storage server
os=$(go env GOOS)   
arch=$(go env GOARCH)
build=
bucket=elabox-debug
#debug host default
echo "OS="$os
echo "Arch="$arch
echo "Upload for version 1 - Staging, 2 - Release, None = Debug"
read answer
# RELEASE
if [ "$answer" == "2" ]; then
    echo "Are you sure to upload release version, this will affect consumer version updates? (y/n)"
    read answer
    if [ "$answer" != "y" ]; then
        exit
    fi
    bucket=elabox
# STAGING
elif [ "$answer" == "1" ]; then
    bucket=elabox-staging
fi

# read build number form system's info.json
pki=../builds/$os/system/info.json
build=$(jq ".build" $pki)

elapath=gs://$bucket
. ~/.bashrc     # reload environment var. there are some instance it is not up to date
gspk=$elapath/packages/$build.box
gspki=$elapath/packages/$build.json
gsinstaller=$elapath/installer/$os/$arch/packageinstaller
gsh=$elapath/installer/$os/$arch/installer.sh

installer=../builds/$os/packageinstaller/bin/packageinstaller
pkg=../builds/$os/system/ela.system.box
shi=./dlinstall.sh
shbk=/tmp/dlinstall.sh

# replace !<variable> from ./dlinstall.sh with dynamic values 
cp -R $shi $shbk
sed -i "s|\!bucket|$bucket|" $shbk
sed -i "s|\!build|$build|" $shbk
sed -i "s|\!arch|$arch|" $shbk

gsutil rm $gsinstaller
gsutil cp $installer $gsinstaller
gsutil cp $pkg $gspk
gsutil cp $pki $gspki
#gsutil cp $pkg $elapath/packages/3.box # remove this later. this is for testing OTA update
#gsutil acl ch -u AllUsers:R $elapath/packages/3.box # remove this later
gsutil cp $shbk $gsh
gsutil acl ch -u AllUsers:R $gspk
gsutil acl ch -u AllUsers:R $gspki
gsutil acl ch -u AllUsers:R $gsinstaller
gsutil acl ch -u AllUsers:R $gsh
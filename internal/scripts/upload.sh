#!/bin/bash
os=linux
arch=arm64
build=2
bucket=elabox-debug
echo "Upload for version 1 - Staging, 2 - Release, None = Debug"
read answer
if [ "$answer" == "2" ]; then
    bucket=elabox
elif [ "$answer" == "1" ]; then
    bucket=elabox-staging
fi
elapath=gs://$bucket

. ~/.bashrc
gspk=$elapath/packages/$build.box
gspki=$elapath/packages/$build.json
gsinstaller=$elapath/installer/$os/$arch/packageinstaller
gsh=$elapath/installer/$os/$arch/installer.sh

installer=../builds/$os/packageinstaller/bin/packageinstaller
pkg=../builds/$os/system/ela.system.box
pki=../builds/$os/system/info.json
shi=./dlinstall.sh
shic=/tmp/ela/dlinstall.sh

cp -R $shi $shic
sed -i "s|\!bucket|$bucket|" $shic

gsutil rm $gsinstaller
gsutil cp $installer $gsinstaller
gsutil cp $pkg $gspk
gsutil cp $pki $gspki
#gsutil cp $pkg $elapath/packages/3.box # remove this later. this is for testing OTA update
#gsutil acl ch -u AllUsers:R $elapath/packages/3.box # remove this later
gsutil cp $shic $gsh
gsutil acl ch -u AllUsers:R $gspk
gsutil acl ch -u AllUsers:R $gspki
gsutil acl ch -u AllUsers:R $gsinstaller
gsutil acl ch -u AllUsers:R $gsh
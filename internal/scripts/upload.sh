#!/bin/bash
os=$(go env GOOS)   
arch=$(go env GOARCH)
build=3
bucket=elabox-debug
#debug host default
rewardhost=208.87.134.80:1235 
echo "OS="$os
echo "Arch="$arch
echo "Upload for version 1 - Staging, 2 - Release, None = Debug"
read answer
if [ "$answer" == "2" ]; then
    bucket=elabox
    rewardhost=208.87.134.80:1234
elif [ "$answer" == "1" ]; then
    bucket=elabox-staging
    rewardhost=208.87.134.80:1236
fi

elapath=gs://$bucket
. ~/.bashrc     # reload environment var. there are some instance it is not up to date
gspk=$elapath/packages/$build.box
gspki=$elapath/packages/$build.json
gsinstaller=$elapath/installer/$os/$arch/packageinstaller
gsh=$elapath/installer/$os/$arch/installer.sh

installer=../builds/$os/packageinstaller/bin/packageinstaller
pkg=../builds/$os/system/ela.system.box
pki=../builds/$os/system/info.json
shi=./dlinstall.sh
shbk=/tmp/dlinstall.sh

# replace !<variable> from ./dlinstall.sh with dynamic values 
cp -R $shi $shbk
sed -i "s|\!bucket|$bucket|" $shbk
sed -i "s|\!build|$build|" $shbk
sed -i "s|\!arch|$arch|" $shbk
sed -i "s|\!rewardhost|$rewardhost|" $shbk

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
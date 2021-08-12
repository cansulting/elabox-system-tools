#!/bin/bash
os=linux
arch=arm64
build=1
. ~/.bashrc
gspk=gs://elabox/packages/$build.box
gspki=gs://elabox/packages/$build.json
gsinstaller=gs://elabox/installer/$os/$arch/packageinstaller
gsh=gs://elabox/installer/$os/$arch/installer.sh

installer=../builds/$os/packageinstaller/bin/packageinstaller
pkg=../builds/$os/system/ela.system.box
pki=../builds/$os/system/info.json
shi=./dlinstall.sh

gsutil rm $gsinstaller
gsutil cp $installer $gsinstaller
gsutil cp $pkg $gspk
gsutil cp $pki $gspki
gsutil cp $pkg gs://elabox/packages/2.box # remove this later. this is for testing OTA update
gsutil acl ch -u AllUsers:R gs://elabox/packages/2.box # remove this later
gsutil cp $shi $gsh
gsutil acl ch -u AllUsers:R $gspk
gsutil acl ch -u AllUsers:R $gspki
gsutil acl ch -u AllUsers:R $gsinstaller
gsutil acl ch -u AllUsers:R $gsh
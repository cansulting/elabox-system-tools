#!/bin/bash
os=linux
arch=arm64
build=1
gspk=gs://elabox/packages/$build.box
gspki=gs://elabox/packages/$build.json
gsinstaller=gs://elabox/installer/$os/$arch/packageinstaller

installer=../builds/$os/packageinstaller/bin/packageinstaller
pkg=../builds/$os/system/ela.system.box
pki=../builds/$os/system/info.json

gsutil rm $gsinstaller
gsutil cp $installer $gsinstaller
gsutil cp $pkg $gspk
gsutil cp $pki $gspki
gsutil acl ch -u AllUsers:R $gspk
gsutil acl ch -u AllUsers:R $gspki
gsutil acl ch -u AllUsers:R $gsinstaller
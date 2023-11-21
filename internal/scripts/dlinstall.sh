#!/bin/bash

build=!build
# select specific build
if [ "$1" != "" ]; then
    build=$1
fi

os=linux
arch=!arch
bucket=!bucket
installer=https://storage.googleapis.com/$bucket/installer/$os/$arch/packageinstaller
pkg=https://storage.googleapis.com/$bucket/packages/$build.box
gen=1  # current generation of device

envpath=/var/ela/data/ela.system/env.json
home=/home/elabox

# reconfigure elabox?
if [[ -f $envpath ]]; then 
    echo "Reconfigure Elabox? (y/n)"
    read resetup
    if [[ $resetup == "y" ]]; then
        sudo rm $envpath
    fi

    # check disk?
    if [[ -d $home ]]; then
        echo "Run check disk for any storage issues? (y/n)"
        read chkdsk

        if [[ $chkdsk == "y" ]]; then
            sudo umount /dev/sda
            sudo fsck /dev/sda
            sudo mount /dev/sda $home
        fi
    fi
fi

echo "Start downloading package"
sudo wget "$pkg"

echo "Start downloading installer..."
sudo wget "$installer" 
sudo chmod +x ./packageinstaller

echo "Installing..."
sudo ./packageinstaller $build.box
# 
echo "Cleaning up..."
sudo rm ./packageinstaller
sudo rm ./$build.box

echo "Rebooting..."
sudo reboot
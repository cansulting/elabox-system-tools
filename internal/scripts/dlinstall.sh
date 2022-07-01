#!/bin/bash

build=!build
# select specific build
if [ "$1" != "" ]; then
    build=$1
fi

os=linux
arch=!arch
installer=https://storage.googleapis.com/!bucket/installer/$os/$arch/packageinstaller
pkg=https://storage.googleapis.com/!bucket/packages/$build.box
gen=1  # current generation of device

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
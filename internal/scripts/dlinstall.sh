#!/bin/bash
os=linux
arch=arm64
build=1
installer=https://storage.googleapis.com/elabox/installer/$os/$arch/packageinstaller
pkg=https://storage.googleapis.com/elabox/packages/$build.box
echo "Start downloading package"
sudo wget "$pkg"

echo "Start downloading installer"
sudo wget "$installer" 
sudo chmod +x ./packageinstaller

echo "Installing"
sudo ./packageinstaller $build.box -r

echo "Delete downloaded (y/n)?"
read delete
if [[ "$delete" == "y" ]]; then
    sudo rm ./packageinstaller
    sudo rm ./$build.box
fi

echo "Reboot? (y/n)"
read rb
if [[ "$rb" == "y" ]]; then 
    sudo reboot
fi
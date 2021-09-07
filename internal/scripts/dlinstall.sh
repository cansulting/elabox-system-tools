#!/bin/bash
os=linux
arch=!arch
build=!build
installer=https://storage.googleapis.com/!bucket/installer/$os/$arch/packageinstaller
pkg=https://storage.googleapis.com/!bucket/packages/$build.box
echo "Start downloading package"
sudo wget "$pkg"

echo "Start downloading installer..."
sudo wget "$installer" 
sudo chmod +x ./packageinstaller

echo "Installing..."
sudo ./packageinstaller $build.box

echo "Cleaning up..."
sudo rm ./packageinstaller
sudo rm ./$build.box

echo "Rebooting..."
sudo reboot
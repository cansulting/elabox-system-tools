#!/bin/bash
# create and compress the specific usb storage as image

shrinkp=/usr/local/bin/pishrink.sh
output=elabox.img
outputz=elabox.zip
bucket=elabox-debug
os=$(go env GOOS)   
arch=$(go env GOARCH)

if [[ ! -f "$shrinkp" ]]; then
    wget https://raw.githubusercontent.com/Drewsif/PiShrink/master/pishrink.sh
    chmod +x pishrink.sh
    sudo mv pishrink.sh /usr/local/bin
fi

echo "ENV select 1 - Staging, 2 - Release, None = Debug"
read env
if [[ "$env" -eq "1" ]]; then 
    bucket=elabox-staging
else 
    if [[ "$env" -eq "2" ]]; then
        bucket=elabox
    fi
fi

echo "Input source storage id for image, use lsblk to view list.(eg: sdb)"
read storage
if [[ $storage -eq "" ]]; then
    storage=sdb
fi

echo "Upload build? (y/n)"
read upload

echo "Please wait, this takes time..."
sudo dd if=/dev/$storage bs=4M conv=sparse of=_$output status=progress
eval "$shrinkp" _$output $output
sudo rm _$output
echo "Image created, @$output..."
du -sch $output

# Upload build
if [[ $upload -eq "y" ]]; then
    echo "Compressing image..."
    sudo zip $outputz $output
    echo "Start uploading..."
    gpath=gs://$bucket/installer/$os/$arch/$outputz
    gsutil cp ./$outputz $gpath
    gsutil acl ch -u AllUsers:R $gpath
    rm $outputz
    echo "Done uploading."
fi
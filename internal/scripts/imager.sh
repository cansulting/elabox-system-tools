#!/bin/bash
# create and compress the specific usb storage

shrinkp=/usr/local/bin/pishrink.sh
output=elabox.img

if [[ ! -f "$shrinkp" ]]; then
    wget https://raw.githubusercontent.com/Drewsif/PiShrink/master/pishrink.sh
    chmod +x pishrink.sh
    sudo mv pishrink.sh /usr/local/bin
fi

echo "Input source storage id for image, use lsblk to view list.(eg: sdb)"
read storage
if [[ $storage -eq "" ]]; then
    storage=sdb
fi

echo "Please wait, this takes time..."
sudo dd if=/dev/$storage bs=4M conv=sparse of=_$output
eval "$shrinkp" _$output $output
sudo rm _$output
echo "Image created, @$output..."
du -sch $output



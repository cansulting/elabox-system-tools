#!/bin/bash
os=linux
arch=!arch
build=!build
installer=https://storage.googleapis.com/!bucket/installer/$os/$arch/packageinstaller
pkg=https://storage.googleapis.com/!bucket/packages/$build.box
rewhost=localhost:1234

# elabox registration. purchase a premium license
for (( ; ; ))
do
    echo ""
    echo "Register your Elabox? (y/n)"
    read license
    if [ "$license" == "y" ]; then
        echo "Input your license:"
        read license
        secret=$license
        gen=1
        serial=$(cat /proc/cpuinfo | grep Serial | cut -d ' ' -f 2)
        hardware=$(cat /proc/cpuinfo | grep Hardware | cut -d ' ' -f 2-10)
        model="$(cat /proc/cpuinfo | grep Model | cut -d ' ' -f 2-10)"
        response=$(curl --location -G \
            --data-urlencode "model=$model" \
            "http://$rewhost/apiv1/rewards/reg-manual?secret=$secret&serial=$serial&hardware=$hardware&gen=$gen" \
            --request POST
        )
        resultCode=$(echo $response | jq '.code')
        if [ "$resultCode" == 200 ]; then
            echo "Registration success!"
            break
        else
            echo "Failed registration ".$response
        fi
    else
        break
    fi
done


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
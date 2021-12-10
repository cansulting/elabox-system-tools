#!/bin/bash
os=linux
arch=!arch
build=!build
installer=https://storage.googleapis.com/!bucket/installer/$os/$arch/packageinstaller
pkg=https://storage.googleapis.com/!bucket/packages/$build.box
rewhost=!rewardhost

# utility package
sudo apt install jq

# check if already registered. return true if registered
isRegistered() {
    serial=$(cat /proc/cpuinfo | grep Serial | cut -d ' ' -f 2)
    response=$(curl --location -G \
            "http://$rewhost/apiv1/rewards/check-device?serial=$serial" \
            --request POST
        )
    resultCode=$(echo $response | jq '.code')
    resultData=$(echo $response | jq '.data')
    if [ "$resultCode" == 200 ]; then
        if [ "$resultData" != "null" ]; then
            echo "true"
        else
            echo false
        fi
    else
        echo "Failed registration ".$response
    fi
}

#########################################
# Registration
#########################################
for (( ; ; ))
do
    echo ""
    checkRes=$(isRegistered)
    # check if already 
    if [ "$checkRes" == "true" ]; then 
        echo "Your elabox was registered."
        break;
    elif [ "$checkRes" != "false" ]; then 
        echo "Failed check." . $checkRes . "Retrying check " 
        continue
    fi

    echo "Register your Elabox? (y/n)"
    read license
    if [ "$license" == "y" ]; then
        echo "Input your license number:"
        read license
        gen=1
        serial=$(cat /proc/cpuinfo | grep Serial | cut -d ' ' -f 2)
        hardware=$(cat /proc/cpuinfo | grep Hardware | cut -d ' ' -f 2-10)
        model="$(cat /proc/cpuinfo | grep Model | cut -d ' ' -f 2-10)"
        response=$(curl --location -G \
            --data-urlencode "model=$model" \
            "http://$rewhost/apiv1/rewards/reg-manual?license=$license&serial=$serial&hardware=$hardware&gen=$gen" \
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
# 
echo "Cleaning up..."
sudo rm ./packageinstaller
sudo rm ./$build.box

echo "Rebooting..."
sudo reboot
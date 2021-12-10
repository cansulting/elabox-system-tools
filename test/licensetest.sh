#!/bin/bash
# use to check if serial was already registered
rewhost=208.87.134.80:1236 #STAGING host

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
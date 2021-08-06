#!/bin/bash
echo "Setting development pipeline for goolang"
echo "Optional commandline params -o(target) -a(arch))"
cos=linux                 
carc=arm64

# CHMOD
sudo chmod +x ./upload.sh
sudo chmod +x ./syncproj.sh
sudo chmod +x ./build.sh
sudo chmod +x ./install.sh

#sudo add-apt-repository ppa:longsleep/golang-backports
#sudo apt update
#sudo apt install golang-go
# FLAGS
while getopts o:a flag
do
    case "${flag}" in
        o) cos=${OPTARG};;
        a) carc=${OPTARG};;
    esac
done

# download go lang
if [ ! -d "/usr/local/go" ]; then 
    wget https://golang.org/dl/go1.16.6.$cos-$carc.tar.gz
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.6.$cos-$carc.tar.gz
    rm go1.16.6.$cos-$carc.tar.gz

    export PATH=$PATH:/usr/local/go/bin
    echo ""export PATH=$PATH:/usr/local/go/bin"" >> ~/.bash_profile

    # install gcc pipelines
    sudo apt install gcc-aarch64-linux-gnu
    sudo apt install gcc-multilib -y
    sudo apt install x86_64-linux-gnu-gcc
    sudo apt-get install gcc-mingw-w64
else
    echo "Golang, GCC libraries installed. skipping..."
fi

# install gcp tools
echo "Do you want to setup environment for package uploading? (y/n)"
read answer
if [[ "$answer" == "y" ]]; then
    sudo apt install python
    cw=$PWD
    cd ~
    echo "Setting up GCP storage for packages"
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-350.0.0-linux-arm.tar.gz
    sudo tar xfz google-cloud-sdk-350.0.0-linux-arm.tar.gz
    ./google-cloud-sdk/install.sh
    ./google-cloud-sdk/bin/gcloud init
    sudo rm google-cloud-sdk-350.0.0-linux-arm.tar.gz
    . ~/.bashrc
    cd $cw

    # for json bash parsing
    sudo apt install jq
fi

######################################
## SETUP GIT PROJECTS
######################################
wd=~
echo "Do you want to setup git projects? (y/n)"
read answer
if [ "$answer" == "y" ]; then
    echo "Set working directory" $(echo $wd) ".(Leave empty if use default.)"
    read answer
    if [ "$answer" != "" ]; then 
        wd=$answer
    fi
    echo "Your git username? "
    read uname
    if [ ! -d "./elabox-companion" ]; then
        cd $wd
        git clone https://$uname@github.com/cansulting/elabox-companion.git
        cd "./elabox-companion"
        git switch Development
    fi 
    if [ ! -d "./landing-page" ]; then
        cd $wd
        git clone https://$uname@github.com/bonhokage06/landing-page.git
        cd "./landing-page"
        git switch main
    fi
fi
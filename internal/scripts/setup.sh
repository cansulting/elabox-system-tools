#!/bin/bash
# script use to setup the project

echo "Setting development pipeline for goolang"
echo "Optional commandline params -o(target) -a(arch))"
cos=$(go env GOOS)                  
carc=$(go env GOARCH) 

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

# NodeJs 
curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
echo 'Y' | sudo apt update 
echo 'Y' | sudo apt install nodejs

# for json bash parsing
if [ "$cos" == "linux" ]; then
    echo 'Y' | sudo apt install jq python zip
elif [ "$cos" == "darwin" ]; then
    brew install jq
    brew install zip
    brew install python
fi

# download go lang
if [ ! -d "/usr/local/go" ]; then 
    pkg=go1.17.5.$cos-$carc.tar.gz
    wget https://golang.org/dl/$pkg
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $pkg
    sudo rm $pkg

    export PATH=$PATH:/usr/local/go/bin
    echo ""export PATH=$PATH:/usr/local/go/bin"" >> ~/.bashrc

    # install gcc pipelines
    echo 'Y' | snap install zig --beta --classic # for cross compiling remove other toolchains
    echo 'Y' | sudo apt install gcc-aarch64-linux-gnu
    echo 'Y' | sudo apt install gcc-multilib -y
    echo 'Y' | sudo apt install gcc-9-x86-64-linux-gnux32 # linux amd64
    echo 'Y' | sudo apt-get install gcc-mingw-w64 # windows
    echo 'Y' | sudo apt-get install gcc-i686-linux-gnu #linux intel
    # for carrier build 
    echo 'Y' | sudo apt-get install build-essential autoconf automake autopoint libtool bison texinfo pkg-config cmake
    . ~/.bashrc
else
    echo "Golang, GCC libraries installed. skipping..."
fi

#########################
# install gcp tools
#########################
echo "Do you want to setup environment for package uploading? (y/n)"
read answer
if [[ "$answer" == "y" ]]; then
    cw=$PWD
    cd ~
    echo "Setting up GCP storage for packages"
    curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-350.0.0-linux-arm.tar.gz
    sudo tar xfz google-cloud-sdk-350.0.0-linux-arm.tar.gz
    sudo ./google-cloud-sdk/install.sh
    sudo ./google-cloud-sdk/bin/gcloud init
    sudo rm google-cloud-sdk-350.0.0-linux-arm.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo ""export PATH=$PATH:/usr/local/go/bin:$PWD/google-cloud-sdk/bin"" >> ~/.bashrc
    . ~/.bashrc
    cd $cw
fi

######################################
## SETUP GIT PROJECTS
######################################
wd=../../../
cd $wd
wd=$PWD

echo "Set working directory" ${wd} ".(Leave empty if use default.)"
read answer
if [ "$answer" != "" ]; then 
    wd=$answer
fi

if [ ! -d "elabox-companion" ]; then
    git clone https://github.com/cansulting/elabox-companion.git
    cd "./elabox-companion"
    git switch Development
fi 
if [ ! -d "elabox-companion-landing" ]; then
    cd $wd
    git clone https://github.com/cansulting/elabox-companion-landing.git
    cd "./elabox-companion-landing"
    git switch main
fi
if [ ! -d "elabox-logs" ]; then
    cd $wd
    git clone https://github.com/cansulting/elabox-logs
    cd elabox-logs
    git switch Development
fi
if [ ! -d "mastodon-hub" ]; then
    cd $wd
    git clone https://github.com/cansulting/mastodon-hub
    cd mastodon-hub
    git switch development
fi
if [ ! -d "Elastos.ELA" ]; then 
    cd $wd
    git clone https://github.com/elastos/Elastos.ELA.git
    cd "./Elastos.ELA"
    git switch master
    /usr/local/go/bin/go mod tidy
fi
if [ ! -d "Elastos.NET.Carrier.Bootstrap" ]; then
    cd $wd
    git clone https://github.com/elastos/Elastos.NET.Carrier.Bootstrap.git
    cd "./Elastos.NET.Carrier.Bootstrap"
    git switch master
fi
if [ ! -d "Elastos.ELA.SideChain.EID" ]; then
    cd $wd
    # added library
    echo "Y" | sudo apt-get install autoconf libudev
    git clone https://github.com/elastos/Elastos.ELA.SideChain.EID.git
fi
if [ ! -d "Elastos.ELA.SideChain.ESC" ]; then
    cd $wd
    # added library
    echo "Y" | sudo apt-get install autoconf libudev
    git clone https://github.com/elastos/Elastos.ELA.SideChain.ESC.git
fi
if [ ! -d "elabox-foundation.lib" ]; then
    cd $wd
    git clone https://github.com/cansulting/elabox-foundation.lib.git
    cd elabox-dapp-store
    git switch development
    npm link
fi
if [ ! -d "elabox-dapp-store" ]; then
    cd $wd
    git clone https://github.com/cansulting/elabox-dapp-store.git
    cd "../elabox-dapp-store"
    git switch development
fi
if [ ! -d "elabox-setup-wizard" ]; then 
    cd $wd
    git clone https://github.com/cansulting/elabox-setup-wizard.git
    cd "../elabox-setup-wizard"
    git switch development
fi
if [ ! -d "elabox-dashboard" ]; then 
    cd $wd
    git clone https://github.com/cansulting/elabox-dashboard.git
    cd "../elabox-dashboard"
    git switch main
fi

######################################
## GLIDE APP
######################################
if [ ! -d "glide-node-server" ]; then
    cd $wd
    git clone https://github.com/glide-finance/glide-node-server.git
fi
if [ ! -d "glide-frontend" ]; then
    cd $wd
    git clone https://github.com/glide-finance/glide-frontend.git
fi

#!/bin/bash
homedir=/home/elabox
user=elabox
passwd=elabox

# terminate running process
echo "Killing running nodes..."
if [ "$(pgrep ela)" != "" ]; then
    sudo kill $(pgrep ela)
fi
sudo pkill geth
sudo pkill ela-bootstrapd

############################
## Setup user, dependent files, caching
############################
exists=$(grep -c "^$user:" /etc/passwd)
if [ "$exists" == 0 ]; then
    echo "Set USB as home? If 'y' please insert USB to your elabox. (y/n)"
    read answer

    echo "Setting up user..."
    echo 'exit' | sudo useradd -p $(openssl passwd -1 $passwd) -m $user
    sudo usermod -aG sudo $user
    sudo -s -u $user
    # download files
    echo "Downloading files"
    echo '$user' | curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -
    echo 'Y' | sudo apt update 
    echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p

    ############################
    ## Firewall
    ############################
    echo "Setting up firewall..."
    # open the different ports with ufw
    sudo ufw default deny incoming
    sudo ufw default allow outgoing
    # SSH port
    sudo ufw allow 22
    # companiont app port
    sudo ufw allow 80
    # elabox back-end port
    sudo ufw allow 3001
    # elabox carrier
    sudo ufw allow 33445
    # ELA DPoS port
    sudo ufw allow 20339
    # ELA port for SPV peers
    sudo ufw allow 20338
    # ELA RPC port
    sudo ufw allow 20336
    # ELA REST port
    sudo ufw allow 20334
    # DID REST port
    sudo ufw allow 20604
    # DID RPC port
    sudo ufw allow 20606
    # DID node port
    sudo ufw allow 20608
    echo 'y' | sudo ufw enable
    
    
    ############################
    ## Setup elabox directory from USB
    ############################
    if [ "$answer" == "y" ]; then
        echo 'y' |  mkfs.ext4 /dev/sda
        sudo mount /dev/sda $homedir
        # check the unique identifier of /dev/sda
        USD_UUID=$(sudo blkid | grep /dev/sda | cut -d '"' -f 2)
        # update the /etc/fstab file to auto-mount the disk on startup
        echo "UUID=${USD_UUID} $homedir ext4 defaults 0 0" | tee -a /etc/fstab > /dev/null
        chown -R elabox:elabox $homedir
        /etc/init.d/avahi-daemon restart

        # carrier directory requirement. delete this later
        mkdir -p /home/elabox/supernode/carrier/
    fi
fi

############################
## Setup packages
############################
isinstalled() {
    pk=$1
    echo $(dpkg-query -W -f='${Status}' $pk 2>/dev/null | grep -c "ok installed")
}
# install function - this install a package if not yet installed
install() {
    if [ $(isinstalled $1) -eq 0 ]; then
        echo 'Y' | sudo apt install $1
    fi
}
install fail2ban
install nodejs
# install and setup tor
if [ $(isinstalled tor) -eq 0 ]; then
    echo 'Y' | sudo apt install tor
    # add the webserver and SSH to tor
    echo ""  | sudo tee -a /etc/tor/torrc
    echo "HiddenServiceDir /var/lib/tor/elabox/"  | sudo tee -a /etc/tor/torrc
    echo "HiddenServicePort 80 127.0.0.1:80" | sudo tee -a /etc/tor/torrc
    echo "HiddenServicePort 22 127.0.0.1:22" | sudo tee -a /etc/tor/torrc
    echo "HiddenServicePort 3001 127.0.0.1:3001" | sudo tee -a /etc/tor/torrc
    echo ""  | sudo tee -a /etc/tor/torrc
    sudo systemctl restart tor@default
fi
# install and setup avahi-daemon
if [ $(isinstalled avahi-daemon) -eq 0 ]; then
    echo 'Y' | sudo apt install avahi-daemon
    # Update hostname to elabox
    echo "Updating hostname..."
    echo "$user" | sudo tee /etc/hostname
    echo "127.0.0.1 $user" | sudo tee /etc/hosts
    sudo hostnamectl set-hostname elabox
    # hostnamectl to check
    /etc/init.d/avahi-daemon restart
    systemctl restart systemd-logind.service
fi

############################
## Setup memory paging 
############################
if [ ! -d "/var/cache/swap" ]; then 
    echo "Setting up cache swap files..."
    sudo mkdir -v /var/cache/swap
    cd /var/cache/swap
    sudo dd if=/dev/zero of=swapfile bs=1K count=4M
    sudo chmod 600 swapfile
    sudo mkswap swapfile
    sudo swapon swapfile
    echo "/var/cache/swap/swapfile none swap sw 0 0" | sudo tee -a /etc/fstab
    top -bn1 | grep -i swap
elif ! grep -q '/var/cache/swap/swapfile none swap sw 0 0' /etc/fstab ; then
    # bug fix for build #2. swapfile was not registered to fstab
    echo "/var/cache/swap/swapfile none swap sw 0 0" | sudo tee -a /etc/fstab
fi

#!/bin/bash
homedir=/home/elabox
user=elabox
passwd=elabox
############################
## Setup user, dependent files, caching
############################
exists=$(grep -c "^$user:" /etc/passwd)
if [ "$exists" == 0 ]; then
    echo "Setting up user..."
    sudo useradd -p $(openssl passwd -1 $passwd) -m $user
    sudo usermod -aG sudo $user
    sudo -s -u $user
    # download files
    echo "Downloading files"
    echo '$user' | curl -sL https://deb.nodesource.com/setup_12.x | sudo -E bash -
    echo 'Y' | sudo apt update && sudo apt install fail2ban avahi-daemon nodejs tor zip
    echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p

    #firewall
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
    # elabox IPC via socket io
    sudo ufw allow 9000
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

    # add the webserver and SSH to tor
    echo ""  | sudo tee -a /etc/tor/torrc
    echo "HiddenServiceDir /var/lib/tor/elabox/"  | sudo tee -a /etc/tor/torrc
    echo "HiddenServicePort 80 127.0.0.1:80" | sudo tee -a /etc/tor/torrc
    echo "HiddenServicePort 22 127.0.0.1:22" | sudo tee -a /etc/tor/torrc
    echo "HiddenServicePort 3001 127.0.0.1:3001" | sudo tee -a /etc/tor/torrc
    echo ""  | sudo tee -a /etc/tor/torrc
    sudo systemctl restart tor@default

    # Update hostname to elabox
    echo "Updating hostname..."
    echo "$user" | sudo tee /etc/hostname
    echo "127.0.0.1 $user" | sudo tee /etc/hosts
    sudo hostnamectl set-hostname elabox
    # hostnamectl to check
    /etc/init.d/avahi-daemon restart
    systemctl restart systemd-logind.service
    
    ############################
    ## Setup elabox directory from USB
    ############################
    if [ ! -d "$homedir" ]; then
        echo "Please insert USB before continuing! Press enter when ready..."
        read answer
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
## Setup setup memory paging and Firewall
############################
if [ ! -d "/var/cache/swap" ]; then 
    echo "Setting up cache swap files..."
    sudo mkdir -v /var/cache/swap
    cd /var/cache/swap
    sudo dd if=/dev/zero of=swapfile bs=1K count=4M
    sudo chmod 600 swapfile
    sudo mkswap swapfile
    sudo swapon swapfile
    top -bn1 | grep -i swap
fi


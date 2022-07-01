#!/bin/bash
homedir=/home/elabox
user=elabox
passwd=elabox

# terminate running process
sudo pkill ela.mainchain
sudo pkill ela.eid
sudo pkill ela.esc
sudo pkill ela-bootstrapd
sudo pkill feedsd

# download files
echo "Setting up nodejs"
echo '$user' | curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
echo 'Y' | sudo apt update 

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
## Setup user, dependent files, caching
############################
exists=$(grep -c "^$user:" /etc/passwd)
if [ "$exists" == 0 ]; then
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
    # Feeds
    sudo ufw allow 10018
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
    # ESC NODE RPC PORT
    sudo ufw allow 20637
    echo 'y' | sudo ufw enable
    
    echo "Setting up user..."
    echo 'exit' | sudo useradd -p $(openssl passwd -1 $passwd) -m $user
    sudo usermod -aG sudo $user
    sudo -s -u $user
    echo fs.inotify.max_user_watches=524288 | sudo tee -a /etc/sysctl.conf && sudo sysctl -p
fi

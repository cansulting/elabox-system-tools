echo "Start uninstalling the system"

echo "Killing processes"
ebox terminate
sudo pkill ela
sudo pkill esc
sudo pkill geth

# data
sudo rm -r /var/ela
#sudo rm -r ~/data
# apps
sudo rm -r /usr/ela
sudo rm -r /home/elabox/apps
# caches
sudo rm -r /tmp/ela
# www fies
sudo rm -r /var/www
# shared librarues
sudo rm -r /usr/local/lib/ela
echo "uninstall success"
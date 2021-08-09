echo "Start uninstalling the system"
# data
sudo rm -r /var/ela
sudo rm -r ~/data
# apps
sudo rm -r /usr/ela
sudo rm -r ~/apps
# caches
sudo rm -r /tmp/ela
# www fies
sudo rm -r /var/www
# shared librarues
sudo rm -r /usr/local/lib/ela

echo "Killing processes"
ebox terminate
sudo kill $(pgrep ela)
sudo kill $(pgrep did)
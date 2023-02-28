root=/usr/local/elabox
sudo ./packageinstaller ./ela.system.box -s
sudo chmod -R 0777 $root
sudo ln -s $root/ela.system/ela.system /usr/local/bin/elabox
echo "******Elabox Installed!*****"
echo "type elabox to run"
echo "setting up ela system as service"
fname=elasystem
fnamex=$fname.service
src=$PWD/$fnamex
srcbin=$PWD/main
target=/lib/systemd/system/$fnamex
r=$(echo $src)
# replace source from file
sed -i "s|\!SOURCE|$srcbin|" $src
# replace username & group
user=$(whoami)
group=$(id -g)
log=/var/log/ela.system
sed -i "s|\!USER|$user|" $src
sed -i "s|\!GROUP|$group|" $src
# replace current working directory
sed -i "s|\!CWD|$PWD|" $src
sed -i "s|\!LOG|$log|" $src

cp $src $target
chmod -x $target
rm $src

systemctl enable $fnamex
systemctl start $fname

# create logs at
#journalctl -f -u $fname
echo "Check system log @" $log

# create symlink elasystem
ln -sf $srcbin /bin/elasystem

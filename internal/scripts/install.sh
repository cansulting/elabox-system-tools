echo "Installing system"
os=$(go env GOOS)
path=../builds/$os/packageinstaller/bin/
echo Running at $path
ebox -t
sudo $path/packageinstaller $path../../system/ela.system.box -s
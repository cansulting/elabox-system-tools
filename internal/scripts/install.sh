echo "Installing system"
os=$(go env GOOS)
path=../builds/$os/packageinstaller/bin/
echo Running at $path
sudo $path/packageinstaller $path../../system/ela.system.box -s
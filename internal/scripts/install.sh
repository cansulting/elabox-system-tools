echo "Installing system"
os=$(go env GOOS)
cd ../builds/$os/bins
echo Running at $PWD
sudo ./packageinstaller ../packager/ela.system.box
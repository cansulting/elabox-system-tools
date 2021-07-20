echo "Setting development pipeline for goolang"
sudo add-apt-repository ppa:longsleep/golang-backports
sudo apt update
sudo apt install golang-go

# install gcc pipelines
sudo apt install gcc-aarch64-linux-gnu
sudo apt install gcc-multilib -y
sudo apt install x86_64-linux-gnu-gcc
sudo apt-get install gcc-mingw-w64
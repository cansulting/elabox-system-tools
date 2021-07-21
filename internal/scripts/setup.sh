echo "Setting development pipeline for goolang"
echo "Optional commandline params -o(target) -a(arch))"
cos=linux                 
carc=arm64
#sudo add-apt-repository ppa:longsleep/golang-backports
#sudo apt update
#sudo apt install golang-go
# FLAGS
while getopts o:a flag
do
    case "${flag}" in
        o) cos=${OPTARG};;
        a) carc=${OPTARG};;
    esac
done
wget https://golang.org/dl/go1.16.6.$cos-$carc.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.6.$cos-$carc.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo ""export PATH=$PATH:/usr/local/go/bin"" >> ~/.bash_profile

# install gcc pipelines
sudo apt install gcc-aarch64-linux-gnu
sudo apt install gcc-multilib -y
sudo apt install x86_64-linux-gnu-gcc
sudo apt-get install gcc-mingw-w64
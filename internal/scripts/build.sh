# current os. 
cos=$(go env GOOS)
# current archi
carc=$(go env GOARCH)
pkg_name=packageinstaller
system_name=system
packager=packager
target=$cos
arch=$carc

# FLAGS
while getopts o:a: flag
do
    case "${flag}" in
        o) target=${OPTARG};;
        a) arch=${OPTARG};;
    esac
done

echo "Optional commandline params -o(target) -a(arch)"
echo "OS="$target
echo "Arch="$arch
echo ""
 
# where binaries will be saved
output=$PWD/../builds/$target/bins
go env -w CGO_ENABLED=1
if [ "$target" == "linux" ]; then
    echo "cgo enabled"
    
fi

# build packager
go env -u GOOS
go env -u GOARCH
echo "Building " $packager
go build -o $output ../cwd/$packager

# build binaries
if [ "$target" == "linux" ]; then
    if [ "$arch" == "arm64" ]; then
        # specific gcc for arm64
        go env -w CC=aarch64-linux-gnu-gcc
    fi
else
    # windows intel
    if [ "$arch" == "386" ]; then
        go env -w CXX=i686-w64-mingw32-g++ 
        go env -w CC=i686-w64-mingw32-gcc
    # windows amd
    else
        go env -w CXX=x86_64-w64-mingw32-g++ 
        go env -w CC=x86_64-w64-mingw32-gcc
    fi
fi
go env -w GOOS=$target 
go env -w GOARCH=$arch
echo "Building " $pkg_name
go build -o $output ../cwd/$pkg_name 
echo "Building " $system_name
go build -o $output ../cwd/$system_name
go env -u CC
go env -u CXX

# build companion app?
build_companion=0
if [ "$build_companion" != "1" ]; then
    echo "Build companion app (y/n)?"
    read answer
    if [ "$answer" == "y" ]; then 
        build_companion=1
    fi
fi 
if [ "$build_companion" == "1" ]; then
    # initialize companion path
    echo "Companion app found @ ELA_COMPANION="$ELA_COMPANION
    if [ "$ELA_COMPANION" == "" ]; then 
        echo "Specify companion app react source directory. Hit enter to not update"
        read new_comp
        if [ "$new_comp" != "" ]; then
            new_comp=$(wslpath $new_comp)
            export ELA_COMPANION=$new_comp
            echo "Path updated."
        fi
    fi
    echo "Start building companion app, please wait this might takes time..." 
    initDir=$PWD
    cd $ELA_COMPANION
    sudo npm install
    sudo npm run build
    cd $initDir
    # move front end and node js back end
    rm -r .../builds/$target/www/companion
    mkdir -p ../builds/$target/www/companion
    mkdir -p ../builds/$target/nodejs/companion
    mv $ELA_COMPANION/build ../builds/$target/www/companion
    cp -r $ELA_COMPANION/src_server ../builds/$target/nodejs/companion
    cp -r $ELA_COMPANION/node_modules ../builds/$target/nodejs/companion
    echo "Build success! Moved to" ../builds/$target/www/companion
    echo "Packaging..."
    pkgerPath=../builds/$cos/bins/$packager
    $pkgerPath ../builds/$target/packager/companion.json
fi 

echo "Start packaging..."
pkgerPath=../builds/$cos/bins/$packager
$pkgerPath ../builds/$target/packager/$pkg_name.json
$pkgerPath ../builds/$target/packager/$system_name.json

go env -u CGO_ENABLED
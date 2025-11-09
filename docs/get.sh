#!/usr/bin/env sh

tag=<tag>
platform=$(uname -s | tr '[:upper:]' '[:lower:]')
echo $platform
case $(uname -m) in
  amd64|x86_64)   arch=amd64 ;;
  arm64|aarch64)  arch=arm64 ;;
esac
echo $arch
url=https://github.com/kbenzie/tuck/releases/download/$tag/tuck-$tag-$platform-$arch.tar.gz
echo $url

curl -o tuck.tar.gz --location $url
tar zxvf tuck.tar.gz tuck
mkdir -p ~/.local/bin
mv tuck ~/.local/bin/tuck
rm tuck.tar.gz

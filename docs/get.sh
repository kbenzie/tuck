#!/usr/bin/env bash

platform=$(uname -s | tr '[:upper:]' '[:lower:]')
arch=$(uname -m)
url=https://github.com/kbenzie/tuck/releases/download/<tag>/tuck-<tag>-$platform-$arch.tar.gz

curl -o tuck.tar.gz --location $url
tar zxvf tuck.tar.gz tuck
mkdir -p ~/.local/bin
mv tuck ~/.local/bin/tuck
rm tuck.tar.gz

#!/bin/bash
set -e

sudo dpkg --add-architecture i386
sudo apt-get install lib32stdc++-7-dev lib32z1-dev libc6-dev-i386 linux-libc-dev:i386 g++-multilib

git clone https://github.com/alliedmodders/ambuild
git clone https://github.com/alliedmodders/metamod-source
git clone https://github.com/alliedmodders/hl2sdk hl2sdk-sdk2013
git clone --recursive https://github.com/alliedmodders/sourcemod

sudo pip install ./ambuild

cp client/shim/sourcehook_hookmangen.h metamod-source/core/sourcehook/

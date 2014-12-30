#!/bin/bash

rm ~/linux*.deb
cd ~/ktsan
make CC='../gcc/install/bin/gcc' -j64 deb-pkg LOCALVERSION=-tsan
cd ~/scripts

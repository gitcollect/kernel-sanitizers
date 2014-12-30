#!/bin/bash

rm ~/linux*.deb
cd ~/ktsan
make CC='../gcc/install/bin/gcc' deb-pkg LOCALVERSION=-tsan
cd ~/scripts

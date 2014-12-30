#!/bin/bash

cd ~/ktsan
make CC='../gcc/install/bin/gcc' -j64 LOCALVERSION=-tsan
cd ~/scripts

#!/bin/bash

cd ~/ktsan
make CC='../gcc/install/bin/gcc' LOCALVERSION=-tsan
cd ~/scripts

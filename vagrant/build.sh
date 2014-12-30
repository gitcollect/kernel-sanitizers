#!/bin/bash

rm ~/linux*.deb
cd ~/ktsan
make -j32 deb-pkg LOCALVERSION=-tsan
cd ~/scripts

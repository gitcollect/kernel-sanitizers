#!/bin/bash

cd ~/ktsan
make -j32 LOCALVERSION=-tsan
cd ~/scripts

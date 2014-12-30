#!/bin/bash

cd ~/ktsan
git diff 8e099d1e8be3f598dcefd04d3cd5eb3673d4e098 > ~/scripts/ktsan.patch
./scripts/checkpatch.pl ~/scripts/ktsan.patch
cd ~/scripts
rm ktsan.patch

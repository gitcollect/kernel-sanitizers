#!/bin/bash

cd ~/vagrant_kasan
vagrant ssh -c "echo 'tsan_run_tests' | sudo tee --append /proc/ktsan_tests > /dev/null"
cd ~/scripts

#!/bin/bash

cd ~/vagrant_kasan
vagrant ssh -c "sudo cat /proc/ktsan_stats"
cd ~/scripts

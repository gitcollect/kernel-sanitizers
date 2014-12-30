#!/bin/bash

ssh -v -i id_rsa -p 10022 -o "StrictHostKeyChecking no" root@localhost "echo 'tsan_run_tests' | tee --append /proc/ktsan_tests > /dev/null"

#!/bin/bash

scp -i ./id_rsa -P 10022 -o "StrictHostKeyChecking no" ./trinity root@localhost:~/

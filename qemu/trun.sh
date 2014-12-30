#!/bin/bash

ssh -i ./id_rsa -p 10022 -o "StrictHostKeyChecking no" root@localhost "trinity --dangerous -q -m -C 16"

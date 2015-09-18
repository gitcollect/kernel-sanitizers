#!/bin/bash

sudo apt-get install easy_install
sudo easy_install pip
sudo pip install virtualenv

virtualenv --no-site-packages sandbox
source ./sandbox/bin/activate
easy_install buildbot-slave

master_host_port=127.0.0.1:9990

buildslave create-slave slave $master_host_port ktsan-slave ktsan
buildslave start slave

deactivate

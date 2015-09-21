#!/bin/bash

sudo apt-get install easy_install
sudo easy_install pip
sudo pip install virtualenv

virtualenv --no-site-packages sandbox
source ./sandbox/bin/activate
pip install buildbot-slave

buildslave create-slave slave 127.0.0.1:9990 ktsan-slave ktsan
buildslave start slave

deactivate

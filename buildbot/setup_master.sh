#!/bin/bash

sudo apt-get install easy_install
sudo easy_install pip
sudo pip install virtualenv

virtualenv --no-site-packages sandbox
source ./sandbox/bin/activate
easy_install buildbot

buildslave create-master master
buildslave start master

deactivate

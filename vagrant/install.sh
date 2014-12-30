#!/bin/bash

green='\e[0;32m'
nocolor='\e[0m'

echo -e "${green}Install: setting up ssh...${nocolor}"

cd ~/vagrant_kasan
vagrant ssh-config > ~/scripts/vagrant-ssh-config
cd ~/scripts

echo -e "${green}Install: set up ssh, removing old kernel deb packages...${nocolor}"

cd ~/vagrant_kasan
vagrant ssh -c "rm ~/linux*.deb"
cd ~/scripts

echo -e "${green}Install: removed old kernel deb packages, now copying new packages...${nocolor}"

#scp -F ./vagrant-ssh-config ~/linux-headers*.deb default:~
#scp -F ./vagrant-ssh-config ~/linux-libc*.deb default:~
scp -F ./vagrant-ssh-config ~/linux-image*tsan_*.deb default:~

echo -e "${green}Install: copied kernel deb packages, now installing...${nocolor}"

cd ~/vagrant_kasan

#vagrant ssh -c "sudo dpkg -i ~/linux-headers*.deb"
#vagrant ssh -c "sudo dpkg -i ~/linux-libc*.deb"
vagrant ssh -c "sudo dpkg -i ~/linux-image*tsan_*.deb"

echo -e "${green}Install: installed kernel deb packages, now reloading vagrant box...${nocolor}"

vagrant reload

cd ~/scripts

echo -e "${green}Install: done.${nocolor}"

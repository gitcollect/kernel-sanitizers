#!/bin/bash

cp qemu/linux.img .
qemu-system-x86_64 -hda linux.img -m 20G -smp 4 -net user,hostfwd=tcp::10022-:22 -net nic -nographic -kernel ../ktsan/arch/x86/boot/bzImage -append "console=ttyS0 root=/dev/sda debug earlyprintk=serial" -enable-kvm 2>&1 | tee /tmp/qemu.log

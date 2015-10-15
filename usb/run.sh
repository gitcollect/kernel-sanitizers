#!/bin/bash

cp overlay.qcow2.bu overlay.qcow2
cp ~/upstream/arch/x86/boot/bzImage ./bzImage
cp ~/upstream/vmlinux ./vmlinux

~/usb/qemu/build/x86_64-softmmu/qemu-system-x86_64 \
	-m 2G \
	-smp 2 \
	-net user,hostfwd=tcp::10021-:22 \
	-net nic \
	-nographic \
	-kernel ./bzImage \
	-append "console=ttyS0 root=/dev/sda debug earlyprintk=serial" \
	-enable-kvm \
	-pidfile vm_pid 2>&1 \
	\
	-hda overlay.qcow2 \
	-hdb ram.qcow2 \
	-serial mon:stdio \
	-device nec-usb-xhci \
	-device usb-redir,chardev=usbchardev,debug=0 \
	-chardev socket,server,id=usbchardev,host=127.0.0.1,port=1336,nodelay,nowait \
	| tee /tmp/qemu.log

#	-kernel ~/upstream/arch/x86/boot/bzImage \
#	-kernel ./bzImage \

#~/usb/qemu/build/x86_64-softmmu/qemu-system-x86_64 --enable-kvm -m 2048 -nographic -hda ~/usb/overlays/overlay_0.qcow2 -device nec-usb-xhci -serial mon:stdio -device usb-redir,chardev=usbchardev,debug=0 -chardev socket,server,id=usbchardev,nowait,path=/tmp/vusbf_0_socket -kernel ~/upstream/arch/x86/boot/bzImage -append "console=ttyS0 root=/dev/sda debug earlyprintk=seria

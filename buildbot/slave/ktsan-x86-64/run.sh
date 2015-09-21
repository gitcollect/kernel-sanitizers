#!/bin/bash

set -eux

QEMU=/usr/libexec/qemu-kvm
GCC_PATH=/home/sanitizers/gcc-5.2.0/install/bin/gcc
IMAGE=/home/sanitizers/qemu/image/wheezy.img
BOOT_TIMEOUT=600

echo @@@BUILD_STEP build@@@
echo

make defconfig
make kvmconfig
cat ../ktsan_config >> .config

make CC=$GCC_PATH LOCALVERSION=-ktsan -j32

COMMIT=$(git log --pretty=format:'%h' -n 1)

echo @@@STEP_TEXT@commit: $COMMIT@@@

echo @@@BUILD_STEP boot@@@
echo

cp -f $IMAGE image.img

rm -f vm.pid
rm -f vm.log

BOOT_START=$(date +%s.%N)

$QEMU \
	-hda image.img \
	-m 20G \
	-smp 4 \
	-net user,hostfwd=tcp::10022-:22 \
	-net nic,vlan=0,model=e1000 \
	-nographic \
	-kernel ./arch/x86/boot/bzImage \
	-append "console=ttyS0 root=/dev/sda debug earlyprintk=serial" \
	-enable-kvm \
	-pidfile vm.pid 2>&1 \
	> vm.log &

sleep 1

# Check if the vm process has started.
kill -0 $(cat vm.pid)

# Ensure that we kill the vm process if something fails.
trap "kill $(cat vm.pid); cat vm.log; exit 1" EXIT

timeout $BOOT_TIMEOUT tail -n +0 -f vm.log | \
	grep -q --line-buffered -m 1 "Starting.*sshd"

BOOT_FINISH=$(date +%s.%N)
BOOT_TIME=$(echo "scale=2; ($BOOT_FINISH - $BOOT_START)/1" | bc)

echo @@@STEP_TEXT@time: $BOOT_TIME@@@

kill $(cat vm.pid)

cat vm.log

trap "" EXIT

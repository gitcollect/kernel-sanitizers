#!/bin/bash

go run fuzzer.go \
	-args="--dangerous -q -m -C 16" \
	-binary=~/qemu/image/wheezy/bin/trinity \
	-bzimage=~/stats/new/bzImage \
	-diskimage=~/scripts/wheezy.img \
	-id_rsa=~/qemu/image/ssh/id_rsa \
	-memmb=4096 \
	-maxcpu=4 \
	-ninst=1 \
	-timeout=15m \
	-port=23505



#	-diskimage=~/qemu/image/wheezy.img \

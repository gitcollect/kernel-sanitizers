#!/bin/bash

go run fuzzer.go \
	-args="--dangerous -q -m -C 16" \
	-binary=/usr/local/google/home/andreyknvl/qemu/image/wheezy/bin/trinity \
	-bzimage=/usr/local/google/home/andreyknvl/stats/new/bzImage \
	-diskimage=/usr/local/google/home/andreyknvl/scripts/wheezy.img \
	-id_rsa=/usr/local/google/home/andreyknvl/qemu/image/ssh/id_rsa \
	-memmb=4096 \
	-maxcpu=4 \
	-ninst=1 \
	-timeout=15m \
	-port=23505



#	-diskimage=/usr/local/google/home/andreyknvl/qemu/image/wheezy.img \

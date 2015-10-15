#!/bin/bash

go run fuzzer.go \
	-args="--dangerous -q -m -C 16" \
	-binary=/usr/local/google/home/andreyknvl/qemu/image/wheezy/bin/trinity \
	-bzimage=/usr/local/google/home/andreyknvl/stats/old/bzImage \
	-diskimage=/usr/local/google/home/andreyknvl/qemu/image/wheezy.img \
	-id_rsa=/usr/local/google/home/andreyknvl/qemu/image/ssh/id_rsa \
	-memmb=20000 \
	-maxcpu=4 \
	-ninst=1 \
	-timeout=15m \
	-port=23506

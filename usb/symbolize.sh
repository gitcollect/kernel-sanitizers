#!/bin/bash

./serial.sh | ~/kernel-sanitizers/kernel_symbolize.py \
		--linux='./' \
		--strip='.../upstream/' \
#		--before=1 --after=1

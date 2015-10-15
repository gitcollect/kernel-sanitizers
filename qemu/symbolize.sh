#!/bin/bash

./serial.sh | ~/kernel-sanitizers/kernel_symbolize.py \
		--linux='~/ktsan' \
		--strip='/.../ktsan/' \
#		--before=1 --after=1

#!/usr/bin/python

import json
import random
import sys

syscalls = set([])

with open('../sys/sys.txt') as f:
  for line in f:
    line = line.strip()
    if line.startswith('#'):
      continue
    if '(' not in line:
      continue
    syscall = line.split('(')[0]
    if ' ' in syscall:
      continue
    if '$' in syscall:
	syscall = syscall.split('$')[0]
    syscalls.add(syscall)

assert len(sys.argv) == 2
n = int(sys.argv[1])

syscalls = list(syscalls)
random.shuffle(syscalls)
syscalls = syscalls[:n]

config = {"enable_syscalls": syscalls}
print json.dumps(config, indent=8)

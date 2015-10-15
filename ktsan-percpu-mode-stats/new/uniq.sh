#!/bin/bash

cat *.log | grep -Eoa "ThreadSanitizer.+$" | sort | uniq

#!/bin/bash

./serial.sh | grep -Eoa "ThreadSanitizer.+$" | sort | uniq

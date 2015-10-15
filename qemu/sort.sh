#!/bin/bash

./serial.sh | grep -Eoa "ThreadSanitizer.+$" | sort | uniq -c | sort -g

#!/bin/bash
# Check disk space usage

echo "Checking disk space..."
df -h / | tail -1 | awk '{print "Root partition usage: " $5}'

# Get usage percentage without %
usage=$(df -h / | tail -1 | awk '{print $5}' | sed 's/%//')

if [ "$usage" -gt 90 ]; then
    echo "WARNING: Disk usage is critical!"
    exit 120
elif [ "$usage" -gt 75 ]; then
    echo "WARNING: Disk usage is high"
    exit 110
else
    echo "Disk space is OK"
    exit 0
fi

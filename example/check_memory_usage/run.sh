#!/bin/bash
# Check memory usage

echo "Checking memory usage..."
free -h | grep Mem | awk '{printf "Memory: %s / %s used (%.0f%%)\n", $3, $2, ($3/$2)*100}'

# Get memory usage percentage
usage=$(free | grep Mem | awk '{printf "%.0f", ($3/$2)*100}')

if [ "$usage" -gt 90 ]; then
    echo "WARNING: Memory usage is critical!"
    exit 1
elif [ "$usage" -gt 75 ]; then
    echo "WARNING: Memory usage is high"
    exit 1
else
    echo "Memory usage is OK"
    exit 0
fi

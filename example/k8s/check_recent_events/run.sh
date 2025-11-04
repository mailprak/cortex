#!/bin/bash
# Check recent Kubernetes events for warnings and errors

NAMESPACE=${NAMESPACE:-default}
MINUTES=${MINUTES:-10}

echo "Checking events in namespace: $NAMESPACE (last $MINUTES minutes)"
echo ""

# Get recent warning events
WARNINGS=$(kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' \
    --field-selector type=Warning 2>/dev/null | tail -n +2)

# Get recent error events
ERRORS=$(kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' \
    -o custom-columns=TIME:.lastTimestamp,TYPE:.type,REASON:.reason,MESSAGE:.message \
    2>/dev/null | grep -i "error\|failed\|backoff" | tail -10)

# Show recent events
echo "Recent events:"
kubectl get events -n "$NAMESPACE" --sort-by='.lastTimestamp' | tail -15

echo ""

# Count critical events
ERROR_COUNT=$(echo "$ERRORS" | grep -v "^$" | wc -l)
WARNING_COUNT=$(echo "$WARNINGS" | grep -v "^$" | wc -l)

if [ "$ERROR_COUNT" -gt 5 ]; then
    echo "❌ Found $ERROR_COUNT critical error events"
    exit 130  # Critical: Many errors
elif [ "$ERROR_COUNT" -gt 0 ]; then
    echo "⚠️  Found $ERROR_COUNT error events"
    exit 120  # Warning: Some errors
elif [ "$WARNING_COUNT" -gt 10 ]; then
    echo "⚠️  Found $WARNING_COUNT warning events"
    exit 110  # Warning: Many warnings
elif [ "$WARNING_COUNT" -gt 0 ]; then
    echo "ℹ️  Found $WARNING_COUNT warning events (acceptable)"
    exit 0
else
    echo "✓ No critical events found"
    exit 0
fi

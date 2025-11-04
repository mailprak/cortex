#!/bin/bash
# Restart a pod (by deleting it, will be recreated by deployment/replicaset)

NAMESPACE=${NAMESPACE:-default}
POD_NAME=${POD_NAME:-$1}

if [ -z "$POD_NAME" ]; then
    echo "❌ Error: POD_NAME environment variable or argument required"
    echo "Usage: POD_NAME=my-pod ./run.sh"
    echo "   or: ./run.sh my-pod"
    exit 1
fi

echo "🔄 Restarting pod: $POD_NAME in namespace: $NAMESPACE"
echo ""

# Check if pod exists
if ! kubectl get pod "$POD_NAME" -n "$NAMESPACE" &>/dev/null; then
    echo "❌ Pod not found: $POD_NAME"
    exit 1
fi

# Get the pod's owner (deployment, statefulset, etc.)
OWNER=$(kubectl get pod "$POD_NAME" -n "$NAMESPACE" -o jsonpath='{.metadata.ownerReferences[0].kind}')

echo "Pod owner: $OWNER"
echo ""

# Delete the pod
if kubectl delete pod "$POD_NAME" -n "$NAMESPACE" --grace-period=30; then
    echo ""
    echo "✓ Pod deleted successfully"

    if [ "$OWNER" != "null" ] && [ -n "$OWNER" ]; then
        echo "ℹ️  Pod will be recreated by $OWNER"
        sleep 5

        # Try to find the new pod
        echo ""
        echo "Checking for new pod..."
        kubectl get pods -n "$NAMESPACE" | grep "${POD_NAME%-*}"
    fi

    exit 0
else
    echo ""
    echo "❌ Failed to delete pod"
    exit 1
fi

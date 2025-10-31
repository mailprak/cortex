#!/bin/bash
# Check Kubernetes Pod Status
# Usage: Set NAMESPACE and POD_NAME env vars or pass as args

NAMESPACE=${NAMESPACE:-default}
POD_NAME=${POD_NAME:-$1}

if [ -z "$POD_NAME" ]; then
    echo "Checking all pods in namespace: $NAMESPACE"
    # Get all pods not in Running state
    NOT_RUNNING=$(kubectl get pods -n "$NAMESPACE" --no-headers | grep -v "Running" | wc -l)
    CRASHLOOP=$(kubectl get pods -n "$NAMESPACE" --no-headers | grep "CrashLoopBackOff" | wc -l)
    PENDING=$(kubectl get pods -n "$NAMESPACE" --no-headers | grep "Pending" | wc -l)

    kubectl get pods -n "$NAMESPACE"

    if [ "$CRASHLOOP" -gt 0 ]; then
        echo ""
        echo "❌ Found $CRASHLOOP pod(s) in CrashLoopBackOff"
        exit 130  # Critical: CrashLoopBackOff
    elif [ "$PENDING" -gt 0 ]; then
        echo ""
        echo "⚠️  Found $PENDING pod(s) in Pending state"
        exit 120  # Warning: Pending pods
    elif [ "$NOT_RUNNING" -gt 0 ]; then
        echo ""
        echo "⚠️  Found $NOT_RUNNING pod(s) not in Running state"
        exit 110  # Warning: Non-running pods
    else
        echo ""
        echo "✓ All pods are Running"
        exit 0
    fi
else
    echo "Checking pod: $POD_NAME in namespace: $NAMESPACE"
    STATUS=$(kubectl get pod "$POD_NAME" -n "$NAMESPACE" -o jsonpath='{.status.phase}' 2>/dev/null)

    if [ -z "$STATUS" ]; then
        echo "❌ Pod not found"
        exit 1
    fi

    kubectl get pod "$POD_NAME" -n "$NAMESPACE"

    case "$STATUS" in
        "Running")
            echo "✓ Pod is Running"
            exit 0
            ;;
        "Pending")
            echo "⚠️  Pod is Pending"
            exit 120
            ;;
        "Failed"|"CrashLoopBackOff")
            echo "❌ Pod is in failed state"
            exit 130
            ;;
        *)
            echo "⚠️  Pod is in $STATUS state"
            exit 110
            ;;
    esac
fi

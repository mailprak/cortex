#!/bin/bash
# Check Kubernetes Node Resources (CPU, Memory, Disk)

echo "Checking node resources..."
echo ""

# Get node resource usage
kubectl top nodes 2>/dev/null || {
    echo "⚠️  metrics-server not available, checking node conditions instead"
    kubectl get nodes -o wide
    exit 110
}

echo ""

# Check for nodes with high resource usage
HIGH_CPU=$(kubectl top nodes --no-headers | awk '{if ($3+0 > 80) print $1}')
HIGH_MEM=$(kubectl top nodes --no-headers | awk '{if ($5+0 > 80) print $1}')

# Check node conditions
NOT_READY=$(kubectl get nodes --no-headers | grep -v " Ready" | wc -l)

if [ "$NOT_READY" -gt 0 ]; then
    echo "❌ Found $NOT_READY node(s) NOT Ready"
    kubectl get nodes | grep -v "Ready"
    exit 130  # Critical: Nodes not ready
fi

if [ -n "$HIGH_CPU" ]; then
    echo "⚠️  Nodes with high CPU (>80%):"
    echo "$HIGH_CPU"
    exit 120  # Warning: High CPU
fi

if [ -n "$HIGH_MEM" ]; then
    echo "⚠️  Nodes with high Memory (>80%):"
    echo "$HIGH_MEM"
    exit 120  # Warning: High Memory
fi

echo "✓ All nodes healthy with normal resource usage"
exit 0

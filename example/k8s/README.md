# Kubernetes Debugging Use Cases with Cortex

This directory contains practical Kubernetes debugging neurons and synapses.

## Available Neurons

### Check Neurons

#### 1. `check_pod_status`
**Purpose:** Check if pods are running properly in a namespace

**Usage:**
```bash
# Check all pods in default namespace
NAMESPACE=default ./check_pod_status/run.sh

# Check specific pod
POD_NAME=my-app-pod NAMESPACE=production ./check_pod_status/run.sh
```

**Exit Codes:**
- `0` - All pods running
- `110` - Some pods not running (warning)
- `120` - Pods in Pending state (warning)
- `130` - Pods in CrashLoopBackOff (critical)

**Use Case:** Monitor application health, detect failing pods

---

#### 2. `check_node_resources`
**Purpose:** Check node CPU and memory usage

**Usage:**
```bash
./check_node_resources/run.sh
```

**Exit Codes:**
- `0` - All nodes healthy
- `110` - metrics-server not available (info)
- `120` - High CPU or memory (>80%) on nodes
- `130` - Nodes not ready (critical)

**Use Case:** Capacity planning, detect resource exhaustion

---

#### 3. `check_recent_events`
**Purpose:** Scan recent Kubernetes events for errors and warnings

**Usage:**
```bash
# Check default namespace
NAMESPACE=default ./check_recent_events/run.sh

# Check production namespace
NAMESPACE=production MINUTES=30 ./check_recent_events/run.sh
```

**Exit Codes:**
- `0` - No critical events
- `110` - Multiple warning events (info)
- `120` - Some error events detected
- `130` - Many critical error events (>5)

**Use Case:** Detect issues early, audit cluster activity

---

### Mutate Neurons

#### 4. `mutate_restart_pod`
**Purpose:** Restart a pod by deleting it (will be recreated by controller)

**Usage:**
```bash
POD_NAME=my-app-pod NAMESPACE=default ./mutate_restart_pod/run.sh
```

**Exit Codes:**
- `0` - Pod restarted successfully
- `1` - Failed to restart pod

**Use Case:** Automatic remediation for CrashLoopBackOff pods

---

## Synapses

### `k8s_cluster_health`
**Purpose:** Comprehensive Kubernetes cluster health check

**What it does:**
1. Checks node resources (CPU, memory)
2. Checks pod status in namespace
3. If pods in CrashLoopBackOff → automatically restarts them
4. Checks recent events for errors

**Usage:**
```bash
cd k8s_cluster_health
../../../cortex exec -p .
```

**Configuration:**
Edit `synapse.yaml` to customize:
- Which namespace to check
- Whether to auto-restart failed pods
- Exit on first error or continue

---

## Real-World Use Cases

### 1. **Application Deployment Smoke Test**
After deploying a new version, run the synapse to ensure:
- All pods are running
- No crash loops
- Nodes have capacity
- No error events

```bash
NAMESPACE=production ../../../cortex exec -p k8s_cluster_health
```

### 2. **Automated Incident Response**
When alerts fire, run diagnostics automatically:
```yaml
# Add to your synapse
plan:
  steps:
    serial:
      - check_node_resources
      - check_pod_status
      - check_recent_events
```

### 3. **Scheduled Health Checks**
Add to cron or Kubernetes CronJob:
```bash
*/5 * * * * cd /path/to/cortex/example/k8s && ./cortex exec -p k8s_cluster_health
```

### 4. **Self-Healing Cluster**
Enable automatic fixes:
```yaml
definition:
  - neuron: check_pod_status
    config:
      fix:
        130: mutate_restart_pod  # Auto-restart crashed pods
```

---

## More Use Case Ideas

### Additional Neurons You Could Create:

1. **check_pvc_status** - Check PersistentVolumeClaims
   - Exit 130: PVC stuck in Pending
   - Exit 120: Low disk space

2. **check_ingress_connectivity** - Test ingress endpoints
   - Exit 130: Ingress unreachable
   - Exit 120: High latency

3. **check_hpa_status** - Check HorizontalPodAutoscaler
   - Exit 130: HPA unable to scale
   - Exit 120: Approaching max replicas

4. **mutate_scale_deployment** - Scale up/down deployments
   - Trigger on high CPU

5. **check_cert_expiry** - Check TLS certificate expiration
   - Exit 120: Cert expires in <30 days
   - Exit 130: Cert expires in <7 days

6. **check_image_pull_errors** - Detect ImagePullBackOff
   - Exit 130: Images can't be pulled

7. **check_network_policies** - Validate NetworkPolicy
   - Test connectivity between services

8. **mutate_cordon_node** - Cordon unhealthy nodes
   - Trigger on node NotReady

---

## Kubernetes-Specific Tips

### Environment Variables
Set these before running neurons:
```bash
export NAMESPACE=production
export KUBECONFIG=/path/to/kubeconfig
```

### RBAC Requirements
The service account needs these permissions:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cortex-debugger
rules:
- apiGroups: [""]
  resources: ["pods", "nodes", "events"]
  verbs: ["get", "list", "delete"]
- apiGroups: ["apps"]
  resources: ["deployments", "replicasets"]
  verbs: ["get", "list"]
- apiGroups: ["metrics.k8s.io"]
  resources: ["nodes", "pods"]
  verbs: ["get", "list"]
```

### Running in CI/CD
```yaml
# .github/workflows/k8s-health-check.yml
- name: K8s Health Check
  run: |
    cd cortex/example/k8s
    ../../cortex exec -p k8s_cluster_health
```

---

## Next Steps

1. Customize the neurons for your environment
2. Add your own neurons for specific app checks
3. Create multiple synapses for different scenarios
4. Integrate with your monitoring/alerting system
5. Run as Kubernetes CronJobs for continuous monitoring

**Pro Tip:** Start with read-only checks (`check_*` neurons) before enabling auto-remediation (`mutate_*` neurons).

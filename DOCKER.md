# Running Cortex in Docker

Cortex can run anywhere using Docker containers. This ensures consistent execution across different environments.

## Quick Start

### 1. Build the Docker Image

```bash
make docker-build
```

Or manually:
```bash
docker build -t cortex:latest .
```

### 2. Run Cortex Commands

```bash
# Show help
make docker-run

# Run specific command
make docker-run ARGS="create-neuron my_check"

# Open interactive shell
make docker-shell
```

## Common Use Cases

### Run K8s Health Check

```bash
# Run K8s cluster health check
make docker-k8s-example

# With custom namespace
docker run --rm -it \
  -v ~/.kube:/root/.kube:ro \
  -e NAMESPACE=production \
  cortex:latest \
  cortex fire-synapse -p /cortex/example/k8s/k8s_cluster_health
```

### Run System Health Check

```bash
make docker-system-example
```

### Interactive Shell

```bash
make docker-shell

# Inside container:
cortex --help
cd /cortex/example/k8s
cortex fire-synapse -p k8s_cluster_health
```

## Docker Compose

### Basic Usage

```bash
# Start cortex container
docker-compose up -d

# Execute commands
docker-compose exec cortex cortex --help
docker-compose exec cortex cortex fire-synapse -p /cortex/example/k8s/k8s_cluster_health

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Scheduled Health Checks

Run health checks every 5 minutes:

```bash
# Start scheduler
docker-compose --profile scheduler up -d

# View scheduler logs
docker-compose logs -f cortex-scheduler
```

### One-Time Job

```bash
# Run health check once
docker-compose --profile job run --rm cortex-job

# With custom namespace
NAMESPACE=production docker-compose --profile job run --rm cortex-job
```

## Volume Mounts

The container mounts:

1. **Kubeconfig**: `~/.kube` → `/root/.kube` (read-only)
   - For K8s access

2. **Examples**: `./example` → `/cortex/example`
   - Pre-built examples

3. **Custom neurons**: `./custom` → `/cortex/custom`
   - Your custom neurons/synapses

## Environment Variables

Set these when running:

```bash
docker run --rm -it \
  -e NAMESPACE=production \
  -e KUBECONFIG=/root/.kube/config \
  -e VERBOSE=2 \
  -v ~/.kube:/root/.kube:ro \
  cortex:latest \
  cortex fire-synapse -p /cortex/example/k8s/k8s_cluster_health
```

## Creating Custom Neurons in Container

```bash
# Start shell
make docker-shell

# Inside container
cd /cortex/custom
cortex create-neuron check_my_app
# Edit the script
vi check_my_app/run.sh
# Update paths in neuron.yaml
vi check_my_app/neuron.yaml
```

## Kubernetes Deployment

Deploy Cortex as a Kubernetes CronJob:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cortex-health-check
spec:
  schedule: "*/5 * * * *"  # Every 5 minutes
  jobTemplate:
    spec:
      template:
        spec:
          serviceAccountName: cortex-sa
          containers:
          - name: cortex
            image: cortex:latest
            command:
              - cortex
              - fire-synapse
              - -p
              - /cortex/example/k8s/k8s_cluster_health
            env:
            - name: NAMESPACE
              value: default
          restartPolicy: OnFailure
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cortex-sa
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cortex-role
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
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cortex-rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: cortex-role
subjects:
- kind: ServiceAccount
  name: cortex-sa
```

Deploy:
```bash
kubectl apply -f k8s-cronjob.yaml
```

## CI/CD Integration

### GitHub Actions

```yaml
name: K8s Health Check
on:
  schedule:
    - cron: '0 */6 * * *'  # Every 6 hours
  workflow_dispatch:

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Build Cortex
        run: make docker-build

      - name: Configure kubectl
        run: |
          echo "${{ secrets.KUBECONFIG }}" > ~/.kube/config

      - name: Run Health Check
        run: make docker-k8s-example
```

### GitLab CI

```yaml
cortex-health-check:
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t cortex:latest .
    - docker run --rm -v ~/.kube:/root/.kube:ro cortex:latest
      cortex fire-synapse -p /cortex/example/k8s/k8s_cluster_health
  only:
    - schedules
```

## Makefile Commands Reference

```bash
make help                 # Show all available commands
make build               # Build binary locally
make docker-build        # Build Docker image
make docker-run          # Run cortex in Docker
make docker-shell        # Open shell in container
make docker-k8s-example  # Run K8s health check
make docker-system-example # Run system health check
make clean               # Clean builds
make install             # Install locally to /usr/local/bin
```

## Troubleshooting

### Kubeconfig Not Found

```bash
# Verify kubeconfig location
ls -la ~/.kube/config

# Run with explicit path
docker run --rm -it \
  -v /path/to/.kube:/root/.kube:ro \
  cortex:latest cortex fire-synapse -p /cortex/example/k8s/k8s_cluster_health
```

### Permission Denied

```bash
# Check RBAC permissions
kubectl auth can-i list pods
kubectl auth can-i delete pods

# Use service account with proper roles
```

### Container Can't Reach K8s API

```bash
# Use host network mode
docker run --rm -it \
  --network host \
  -v ~/.kube:/root/.kube:ro \
  cortex:latest cortex fire-synapse -p /cortex/example/k8s/k8s_cluster_health
```

## Best Practices

1. **Use read-only mounts** for kubeconfig (`-v ~/.kube:/root/.kube:ro`)
2. **Set resource limits** in production
3. **Use specific image tags**, not `:latest` in production
4. **Run as non-root user** (modify Dockerfile if needed)
5. **Mount custom neurons** via volumes, don't rebuild image
6. **Use secrets** for sensitive data, not environment variables
7. **Test in staging** before enabling auto-remediation

## Running on Different Platforms

### AWS ECS/Fargate
- Build and push to ECR
- Create task definition with cortex image
- Schedule as cron task

### Azure Container Instances
```bash
az container create \
  --resource-group mygroup \
  --name cortex \
  --image cortex:latest \
  --restart-policy Never
```

### Google Cloud Run
```bash
gcloud run jobs create cortex-job \
  --image cortex:latest \
  --schedule="0 */6 * * *"
```

## Next Steps

1. Build the image: `make docker-build`
2. Test locally: `make docker-k8s-example`
3. Deploy to K8s as CronJob
4. Monitor logs and refine neurons
5. Add your custom neurons to `/cortex/custom`

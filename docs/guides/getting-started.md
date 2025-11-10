# Getting Started with Cortex

Quick start guide to get Cortex up and running in 5 minutes.

## What is Cortex?

Cortex is an infrastructure debugging orchestrator that helps you:
- Organize debugging scripts into reusable "neurons"
- Chain neurons into automated "synapses" (workflows)
- Generate debugging scripts from natural language using AI
- Visualize execution with a modern web UI

## Installation

### Option 1: Download Binary (Recommended)

```bash
# Linux/macOS
curl -LO https://github.com/anoop2811/cortex/releases/latest/cortex
chmod +x cortex
sudo mv cortex /usr/local/bin/

# Verify installation
cortex --version
```

### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/anoop2811/cortex.git
cd cortex

# Build binary
make build

# The binary is now at ./cortex
./cortex --version
```

### Option 3: Docker

```bash
docker pull ghcr.io/anoop2811/cortex:latest
docker run --rm cortex --help
```

## Your First Neuron

A **neuron** is a discrete debugging task (shell script + config).

### 1. Create a Neuron

```bash
cortex create-neuron check-nginx
```

This creates:
```
check-nginx/
├── config.yml       # Neuron configuration
└── run.sh          # Execution script
```

### 2. Edit the Neuron

```yaml
# check-nginx/config.yml
---
name: check_nginx
type: check
description: "Check if nginx is running"
exec_file: run.sh
pre_exec_debug: "Checking nginx status..."
post_exec_success_debug: "Nginx is running!"
post_exec_fail_debug:
  1: "Nginx is not running"
```

```bash
# check-nginx/run.sh
#!/bin/bash
systemctl is-active nginx
```

### 3. Execute the Neuron

```bash
cortex exec -p check-nginx
```

Output:
```
Checking nginx status...
✓ Nginx is running!
Exit code: 0
```

## Your First Synapse

A **synapse** is a workflow that chains multiple neurons together.

### 1. Create a Synapse

```bash
cortex create-synapse health-check
```

This creates:
```
health-check/
├── config.yml       # Synapse configuration
└── neurons/         # Neurons directory
```

### 2. Configure the Synapse

```yaml
# health-check/config.yml
---
name: health-check
description: "Complete system health check"
neurons:
  - check-nginx
  - check-database
  - check-disk-space
execution: sequential  # Run neurons in order
stopOnError: true      # Stop if any neuron fails
```

### 3. Execute the Synapse

```bash
cortex exec -p health-check
```

## AI-Powered Neuron Generation

Generate neurons from natural language (requires API key).

### 1. Configure AI Provider

```bash
# OpenAI
export OPENAI_API_KEY="your-key-here"

# Or Anthropic
export ANTHROPIC_API_KEY="your-key-here"

# Or use Ollama (local, no API key needed)
ollama serve
```

### 2. Generate a Neuron

```bash
cortex generate-neuron \
  --prompt "Check if port 8080 is open and listening" \
  --provider openai \
  --output check-port-8080
```

Cortex will:
1. Analyze your prompt
2. Generate a working shell script
3. Create neuron config
4. Validate the syntax
5. Save to `check-port-8080/`

### 3. Review and Execute

```bash
# Review generated code
cat check-port-8080/run.sh

# Execute
cortex exec -p check-port-8080
```

## Web UI Dashboard

Launch the web interface for real-time monitoring.

```bash
cortex ui --port 8080
```

Then open http://localhost:8080 in your browser.

Features:
- **Dashboard**: See all neurons and synapses
- **Real-time logs**: Watch execution in real-time via WebSocket
- **Visual builder**: Drag-and-drop synapse creation
- **Fleet management**: Monitor multiple edge devices

## Common Workflows

### Daily Health Check

```bash
# Create synapse with multiple checks
cortex create-synapse daily-checks

# Configure neurons
cat > daily-checks/config.yml <<EOF
name: daily-checks
neurons:
  - check-nginx
  - check-database
  - check-disk-space
  - check-memory
  - check-cpu
execution: parallel  # Run all at once
EOF

# Execute daily
cortex exec -p daily-checks
```

### Kubernetes Debugging

```bash
# Generate K8s debugging neurons
cortex generate-neuron \
  --prompt "Check all pods in default namespace" \
  --provider ollama

cortex generate-neuron \
  --prompt "Get logs from failed pods" \
  --provider ollama

# Create synapse
cortex create-synapse k8s-debug
```

### Automated Remediation

```yaml
# remediation-synapse/config.yml
name: auto-fix-nginx
neurons:
  - name: check-nginx
  - name: restart-nginx
    condition: "previous.exitCode != 0"  # Only if check failed
  - name: verify-nginx
    condition: "previous.exitCode == 0"  # Only if restart succeeded
```

## Configuration

### Global Config

```bash
# Create config file
mkdir -p ~/.cortex
cat > ~/.cortex/config.yml <<EOF
---
ai:
  default_provider: ollama
  providers:
    openai:
      api_key: "${OPENAI_API_KEY}"
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
    ollama:
      endpoint: "http://localhost:11434"

ui:
  port: 8080
  enable_websocket: true

execution:
  default_timeout: 300s
  max_concurrent: 5
EOF
```

### Environment Variables

```bash
export CORTEX_CONFIG="~/.cortex/config.yml"
export CORTEX_LOG_LEVEL="info"  # debug, info, warn, error
export CORTEX_NEURON_PATH="./neurons"
export CORTEX_SYNAPSE_PATH="./synapses"
```

## Next Steps

- **[User Guide](user-guide.md)** - Complete feature documentation
- **[Contributing](contributing.md)** - Contribute to Cortex
- **[Architecture](../architecture/overview.md)** - Understand how Cortex works
- **[Examples](../../examples/)** - More example neurons and synapses

## Troubleshooting

### Binary not found
```bash
# Make sure binary is in PATH
which cortex

# Or use full path
/usr/local/bin/cortex --help
```

### Permission denied
```bash
chmod +x cortex
```

### AI generation fails
```bash
# Check API key
echo $OPENAI_API_KEY

# Test Ollama connection
curl http://localhost:11434/api/tags

# Use --debug flag
cortex generate-neuron --prompt "..." --debug
```

### Tests failing
```bash
# Install test dependencies
make install-deps

# Run tests
make test-all
```

## Getting Help

- **GitHub Issues**: https://github.com/anoop2811/cortex/issues
- **Discussions**: https://github.com/anoop2811/cortex/discussions
- **Examples**: See `examples/` directory

---

**Ready to dive deeper?** Check out the [User Guide](user-guide.md) for complete documentation.

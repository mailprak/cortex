# Cortex Quick Start Guide

This guide will help you get started with Cortex, an infrastructure debug orchestrator.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/anoop2811/cortex
cd cortex
```

2. Build the binary:
```bash
go build -o cortex .
```

3. (Optional) Move to your PATH:
```bash
sudo mv cortex /usr/local/bin/
```

## Basic Concepts

- **Neuron**: A discrete debugging task (check or mutate operation)
- **Synapse**: An orchestration plan that executes multiple neurons in serial or parallel

## Quick Example

### 1. Create a neuron (check operation)

```bash
./cortex create-neuron check_disk_space
```

This creates a folder with:
- `neuron.yaml` - Configuration file
- `run.sh` - Script to execute (Linux/Mac)
- `run.ps1` - Script to execute (Windows)

### 2. Edit the neuron script

Edit `check_disk_space/run.sh` to add your check logic:

```bash
#!/bin/bash
echo "Checking disk space..."
df -h / | tail -1 | awk '{print "Root partition usage: " $5}'

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
```

### 3. Update the neuron configuration

Edit `check_disk_space/neuron.yaml` to specify allowed exit codes and update the exec_file path:

```yaml
name: check_disk_space
type: check
description: "Check disk space usage"
exec_file: /absolute/path/to/check_disk_space/run.sh
pre_exec_debug: "Checking disk space..."
assert_exit_status:
  - 0
  - 110  # Warning but acceptable
post_exec_success_debug: "Disk space check completed"
post_exec_fail_debug:
  120: "Critical disk usage detected"
```

### 4. Create a synapse

```bash
./cortex create-synapse system_health_check
```

### 5. Configure the synapse

Edit `system_health_check/synapse.yaml`:

```yaml
name: system_health_check
definition:
  - neuron: check_disk_space
    config:
      path: /path/to/check_disk_space
  - neuron: check_memory_usage
    config:
      path: /path/to/check_memory_usage
plan:
  config:
    exit_on_first_error: false
  steps:
    serial:
      - check_disk_space
      - check_memory_usage
    parallel: []
```

### 6. Fire the synapse

```bash
./cortex fire-synapse -p system_health_check
```

## Available Commands

- `cortex create-neuron <name>` - Create a new neuron
  - Use `-t mutate` flag for mutating neurons
- `cortex create-synapse <name>` - Create a new synapse
- `cortex fire-synapse -p <path>` - Execute a synapse
  - Use `-v <level>` for verbose output (0-4)

## Example Workflow

See the `example/` directory for a complete working example with:
- `check_disk_space/` - Disk space check neuron
- `check_memory_usage/` - Memory usage check neuron
- `system_health_check/` - Synapse that orchestrates both checks

To run the example:

```bash
cd example/system_health_check
../../cortex fire-synapse -p .
```

## Exit Codes

Neurons communicate results via exit codes:
- `0` - Success
- Non-zero - Error or warning (define in `assert_exit_status` to treat as acceptable)

You can trigger automatic fixes based on exit codes in the synapse configuration:

```yaml
definition:
  - neuron: check_disk_space
    config:
      path: /path/to/check_disk_space
      fix:
        120: cleanup_disk  # Run cleanup_disk neuron if exit code is 120
```

## Tips

1. Use `check_` prefix for read-only neurons
2. Use `mutate_` prefix for neurons that change state
3. Use absolute paths in neuron.yaml `exec_file` field
4. Test neurons individually before adding to synapse
5. Use `exit_on_first_error: true` for critical checks

## Next Steps

- Review the [README.md](README.md) for detailed architecture
- Check out the example neurons in `example/`
- Create your own neurons for your infrastructure debugging needs

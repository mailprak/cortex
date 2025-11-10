# Cortex Architecture Overview

High-level architecture and design principles for the Cortex infrastructure debugging orchestrator.

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Cortex System                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  ┌──────────────┐                    ┌──────────────────┐       │
│  │ CLI Interface│                    │  Web UI          │       │
│  │  (Cobra)     │                    │  (React + Vite)  │       │
│  └──────┬───────┘                    └────────┬─────────┘       │
│         │                                     │                 │
│         └──────────────┬──────────────────────┘                 │
│                        │                                        │
│              ┌─────────▼──────────┐                            │
│              │   Core Engine       │                            │
│              │  ┌────────────────┐ │                            │
│              │  │ Neuron Exec    │ │                            │
│              │  │ (Shell runner) │ │                            │
│              │  └────────────────┘ │                            │
│              │  ┌────────────────┐ │                            │
│              │  │ Synapse DAG    │ │                            │
│              │  │ (Orchestrator) │ │                            │
│              │  └────────────────┘ │                            │
│              │  ┌────────────────┐ │                            │
│              │  │ AI Generator   │ │ (Future)                   │
│              │  │ (LLM Provider) │ │                            │
│              │  └────────────────┘ │                            │
│              └────────┬────────────┘                            │
│                       │                                          │
│      ┌────────────────┼────────────────┐                        │
│      │                │                │                        │
│  ┌───▼──────┐    ┌───▼─────┐    ┌────▼──────┐                 │
│  │ Database │    │ Storage │    │AI Providers│                 │
│  │ Layer    │    │ Layer   │    │ (Optional) │                 │
│  │          │    │         │    │            │                 │
│  │SQLite or │    │Local or │    │OpenAI,     │                 │
│  │Postgres  │    │S3       │    │Anthropic,  │                 │
│  │          │    │         │    │Ollama      │                 │
│  └──────────┘    └─────────┘    └────────────┘                 │
└─────────────────────────────────────────────────────────────────┘
```

## Core Concepts

### 1. Neurons

**What:** A neuron is a discrete debugging task - a shell script with metadata.

**Structure:**
```
neuron/
├── config.yml       # Neuron metadata
└── run.sh          # Execution script
```

**Example config.yml:**
```yaml
name: check_nginx
type: check
description: "Check if nginx is running"
exec_file: run.sh
pre_exec_debug: "Checking nginx status..."
assertExitStatus: [0]
post_exec_success_debug: "Nginx is running!"
post_exec_fail_debug:
  1: "Nginx is not running"
```

**Responsibilities:**
- Execute shell command
- Validate exit codes
- Provide debug output
- Handle errors gracefully

### 2. Synapses

**What:** A synapse is a workflow that orchestrates multiple neurons in a DAG.

**Structure:**
```
synapse/
├── config.yml       # Workflow definition
└── neurons/         # Referenced neurons
```

**Example config.yml:**
```yaml
name: health-check
description: "System health check workflow"
neurons:
  - check-nginx
  - check-database
  - check-disk-space
execution: sequential  # or parallel
stopOnError: true
```

**Responsibilities:**
- Parse DAG dependencies
- Execute neurons in correct order
- Handle parallel execution
- Manage error propagation
- Track execution state

### 3. AI Generation (Future)

**What:** Generate neurons from natural language using LLMs.

**Supported Providers:**
- OpenAI (GPT-4)
- Anthropic (Claude)
- Ollama (Local)
- Azure OpenAI

**Process:**
1. Parse user prompt
2. Gather context (existing neurons)
3. Call LLM with engineered prompt
4. Generate shell script + config
5. Validate syntax
6. Save neuron

See [AI Neuron Generation](ai-neuron-generation.md) for details.

### 4. Web UI (Future)

**What:** Modern React dashboard for visualization and management.

**Features:**
- Real-time execution logs (WebSocket)
- Visual synapse builder (drag-drop)
- Neuron library browser
- Fleet management (edge devices)
- Mobile-responsive PWA

See [Web UI Architecture](web-ui.md) for details.

## Technology Stack

### Backend (Go)

```
internal/
├── neuron/          # Neuron execution engine
│   ├── neuron.go
│   └── executor.go
├── synapse/         # Synapse orchestration
│   ├── synapse.go
│   ├── dag.go
│   └── executor.go
├── ai/              # AI generation (future)
│   ├── generator.go
│   ├── providers/
│   │   ├── openai.go
│   │   ├── anthropic.go
│   │   └── ollama.go
│   └── prompt.go
├── api/             # REST API (future)
│   ├── handlers.go
│   └── middleware.go
└── db/              # Database layer (future)
    ├── neuron.go
    └── synapse.go
```

**Key Libraries:**
- **Cobra** - CLI framework
- **Viper** - Configuration
- **Zerolog** - Structured logging
- **GORM** - ORM (future)
- **Chi Router** - HTTP routing (future)

### Frontend (Future)

```
web/
├── src/
│   ├── components/      # React components
│   ├── pages/           # Page components
│   ├── hooks/           # Custom hooks
│   ├── api/             # API client
│   └── App.tsx
├── vite.config.ts
└── package.json
```

**Key Libraries:**
- **React 19** - UI framework
- **Vite 6** - Build tool
- **TanStack Query** - Server state
- **Zustand** - Client state
- **React Flow** - DAG visualization
- **Tailwind CSS** - Styling

## Data Flow

### Neuron Execution

```
User Command
    │
    ▼
CLI Parser (Cobra)
    │
    ▼
Neuron Loader
    │
    ▼
Validation
    │
    ▼
Shell Executor
    │
    ├─→ pre_exec_debug (stdout)
    ├─→ run.sh execution
    ├─→ exit code check
    └─→ post_exec_*_debug (stdout)
    │
    ▼
Result (exit code + output)
```

### Synapse Execution

```
User Command
    │
    ▼
CLI Parser
    │
    ▼
Synapse Loader
    │
    ▼
DAG Builder
    │
    ├─→ Parse dependencies
    ├─→ Topological sort
    └─→ Execution plan
    │
    ▼
Executor
    │
    ├─→ Sequential: Run one by one
    │   └─→ Stop on error (optional)
    │
    └─→ Parallel: Run concurrent
        └─→ Wait for all / first failure
    │
    ▼
Aggregate Results
```

### AI Generation (Future)

```
User Prompt
    │
    ▼
Prompt Analyzer
    │
    ├─→ Extract intent
    ├─→ Gather context
    └─→ Estimate cost
    │
    ▼
LLM Provider
    │
    ├─→ OpenAI: gpt-4
    ├─→ Anthropic: claude-3
    └─→ Ollama: local model
    │
    ▼
Response Parser
    │
    ├─→ Extract shell script
    ├─→ Generate config.yml
    └─→ Validate syntax
    │
    ▼
File Writer
    │
    └─→ Save neuron/
        ├── config.yml
        └── run.sh
```

## Deployment Models

### 1. Single Binary (Local)

```bash
# Build
go build -o cortex .

# Run
./cortex exec -p my-neuron
```

**Use Cases:**
- Local development
- Edge devices (Raspberry Pi)
- CI/CD pipelines
- Air-gapped environments

### 2. Kubernetes (Clustered)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cortex
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: cortex
        image: cortex:latest
        ports:
        - containerPort: 8080
```

**Use Cases:**
- Production deployments
- Multi-team environments
- High availability
- Centralized management

## Design Principles

### 1. **Edge-First**

- Single binary (50MB)
- Minimal dependencies
- Low resource usage (256MB RAM)
- Works offline

### 2. **Shell-Native**

- Embrace existing scripts
- No proprietary formats
- Easy migration
- Language-agnostic

### 3. **Progressive Enhancement**

- Core: CLI neuron execution
- Layer 1: Synapse orchestration
- Layer 2: AI generation
- Layer 3: Web UI
- Layer 4: Fleet management

### 4. **Privacy-First**

- AI is optional
- User controls data
- Local-only mode (Ollama)
- No telemetry by default

### 5. **Test-Driven**

- Outer loop: Acceptance tests
- Inner loop: Unit tests
- 90% coverage minimum
- Continuous testing

## Security Considerations

### Command Injection

**Risk:** User input in shell commands
**Mitigation:**
- Validate all inputs
- Escape special characters
- Use allowlists for commands
- Reject suspicious patterns

### API Key Management

**Risk:** Exposed AI provider keys
**Mitigation:**
- Environment variables only
- Never log keys
- Rotate regularly
- Use key vaults in production

### Execution Isolation

**Risk:** Malicious neuron scripts
**Mitigation:**
- Sandboxed execution (future)
- Resource limits (CPU, memory)
- Timeout enforcement
- Permission checks

## Performance Targets

| Metric | Target | Current |
|--------|--------|---------|
| Neuron execution overhead | < 50ms | TBD |
| Synapse DAG build time | < 100ms | TBD |
| AI generation time | < 5s | TBD |
| Web UI load time | < 2s | TBD |
| WebSocket latency | < 100ms | TBD |

## Observability

### Logging (Zerolog)

```go
log.Info().
    Str("neuron", name).
    Int("exitCode", code).
    Msg("Neuron executed")
```

### Metrics (Prometheus) - Future

- Neuron execution count
- Success/failure rate
- Execution duration (histogram)
- AI generation cost

### Tracing (OpenTelemetry) - Future

- Request tracing
- Distributed tracing
- Performance profiling

## Extensibility

### Plugin System (Future)

```
plugins/
├── providers/       # AI provider plugins
├── validators/      # Custom validators
├── reporters/       # Output formatters
└── hooks/          # Lifecycle hooks
```

### Custom Neurons

Users can create custom neuron types:
- Health checks
- Mutations (deployments)
- Diagnostics
- Remediations

### Integration Points

- Webhooks (synapse events)
- API endpoints
- Custom commands
- External storage

## Future Roadmap

### Phase 1: MVP (Current)
- ✅ Core neuron execution
- ✅ Synapse orchestration
- ✅ CLI interface

### Phase 2: AI & UI (Next 6 months)
- 🚧 AI neuron generation
- 🚧 Web UI dashboard
- 🚧 Visual synapse builder

### Phase 3: Scale (6-12 months)
- 📋 Fleet management
- 📋 Plugin marketplace
- 📋 Self-healing mode

### Phase 4: Enterprise (12+ months)
- 📋 Multi-tenancy
- 📋 RBAC
- 📋 Audit logging
- 📋 SaaS offering

---

For detailed specifications, see:
- [AI Neuron Generation](ai-neuron-generation.md)
- [Web UI Architecture](web-ui.md)

For development, see:
- [Contributing Guide](../guides/contributing.md)
- [Testing Guide](../TESTING.md)

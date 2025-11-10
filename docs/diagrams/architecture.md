# Cortex Architecture Diagrams

Visual representations of Cortex's architecture and workflows.

## System Architecture

```mermaid
graph TB
    subgraph "User Interfaces"
        CLI[CLI Interface<br/>cobra + viper]
        WebUI[Web Dashboard<br/>React + Vite]
    end

    subgraph "Core Engine"
        Orchestrator[Orchestrator]
        NeuronExec[Neuron Executor]
        SynapseDAG[Synapse DAG Engine]
        AIGen[AI Generator<br/>OpenAI/Anthropic/Ollama]
    end

    subgraph "Storage Layer"
        FileSystem[(File System<br/>YAML + Shell)]
        Database[(SQLite/Postgres<br/>Future)]
    end

    subgraph "External Services"
        LLM[LLM Providers<br/>OpenAI/Anthropic/Ollama]
        Git[Git/GitHub]
    end

    CLI --> Orchestrator
    WebUI --> Orchestrator
    Orchestrator --> NeuronExec
    Orchestrator --> SynapseDAG
    Orchestrator --> AIGen

    NeuronExec --> FileSystem
    SynapseDAG --> FileSystem
    AIGen --> LLM
    AIGen --> FileSystem

    NeuronExec -.->|future| Database
    SynapseDAG -.->|future| Database

    style CLI fill:#e1f5ff,color:#1f2937
    style WebUI fill:#e1f5ff,color:#1f2937
    style Orchestrator fill:#fff4e1,color:#1f2937
    style NeuronExec fill:#e8f5e9,color:#1f2937
    style SynapseDAG fill:#e8f5e9,color:#1f2937
    style AIGen fill:#f3e5f5,color:#1f2937
    style FileSystem fill:#fce4ec,color:#1f2937
    style LLM fill:#fff3e0,color:#1f2937
```

## Neuron Execution Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Executor
    participant Neuron
    participant Shell

    User->>CLI: cortex exec -p my-neuron
    CLI->>Executor: Load neuron config
    Executor->>Neuron: Read config.yml
    Neuron-->>Executor: Config loaded

    Executor->>Shell: Execute pre_exec_debug
    Shell-->>User: "Checking nginx status..."

    Executor->>Shell: Execute run.sh
    Shell->>Shell: systemctl is-active nginx
    Shell-->>Executor: Exit code: 0

    alt Exit code matches assertExitStatus
        Executor->>Shell: Execute post_exec_success_debug
        Shell-->>User: "✓ Nginx is running!"
    else Exit code doesn't match
        Executor->>Shell: Execute post_exec_fail_debug[exit_code]
        Shell-->>User: "✗ Nginx is not running"
    end

    Executor-->>CLI: Execution complete
    CLI-->>User: Exit with neuron's exit code
```

## Synapse Workflow (Sequential)

```mermaid
flowchart TD
    Start([User executes synapse]) --> LoadConfig[Load synapse config.yml]
    LoadConfig --> ParseNeurons[Parse neuron list]
    ParseNeurons --> BuildDAG[Build dependency graph]

    BuildDAG --> Sequential{Execution mode?}

    Sequential -->|sequential| N1[Execute Neuron 1<br/>check-nginx]
    N1 --> Check1{Success?}

    Check1 -->|Yes| N2[Execute Neuron 2<br/>check-database]
    Check1 -->|No + stopOnError| Fail1[Report failure]

    N2 --> Check2{Success?}
    Check2 -->|Yes| N3[Execute Neuron 3<br/>check-disk-space]
    Check2 -->|No + stopOnError| Fail2[Report failure]

    N3 --> Check3{Success?}
    Check3 -->|Yes| Success[All neurons passed]
    Check3 -->|No| Fail3[Report failure]

    Sequential -->|parallel| Parallel[Execute all neurons<br/>in parallel]
    Parallel --> WaitAll[Wait for all to complete]
    WaitAll --> AllSuccess{All succeeded?}
    AllSuccess -->|Yes| Success
    AllSuccess -->|No| FailParallel[Report failures]

    Success --> End([Exit 0])
    Fail1 --> End1([Exit with failure code])
    Fail2 --> End1
    Fail3 --> End1
    FailParallel --> End1

    style Start fill:#e1f5ff,color:#1f2937
    style Success fill:#c8e6c9,color:#1f2937
    style Fail1 fill:#ffcdd2,color:#1f2937
    style Fail2 fill:#ffcdd2,color:#1f2937
    style Fail3 fill:#ffcdd2,color:#1f2937
    style FailParallel fill:#ffcdd2,color:#1f2937
    style End fill:#c8e6c9,color:#1f2937
    style End1 fill:#ffcdd2,color:#1f2937
```

## AI Neuron Generation Flow

```mermaid
flowchart LR
    subgraph "Input"
        User[User Prompt<br/>'Find process using port 8080']
    end

    subgraph "Analysis"
        Intent[Intent Analysis<br/>- Extract requirements<br/>- Classify type: check/mutate<br/>- Identify technologies]
        Context[Context Gathering<br/>- Load similar neurons<br/>- Check templates<br/>- System info]
    end

    subgraph "Generation"
        Prompt[Build Prompt<br/>- System instructions<br/>- Few-shot examples<br/>- Context]
        LLM[LLM Provider<br/>OpenAI/Anthropic/Ollama]
        Parse[Parse Response<br/>- Extract YAML<br/>- Extract shell script<br/>- Extract README]
    end

    subgraph "Validation"
        ValidateYAML[YAML Schema Check]
        ValidateShell[Shell Syntax Check<br/>shellcheck]
        Security[Security Scan<br/>- No hardcoded secrets<br/>- No dangerous commands]
    end

    subgraph "Output"
        Files[Generated Files<br/>- config.yml<br/>- run.sh<br/>- README.md<br/>- tests]
    end

    User --> Intent
    Intent --> Context
    Context --> Prompt
    Prompt --> LLM
    LLM --> Parse
    Parse --> ValidateYAML
    ValidateYAML --> ValidateShell
    ValidateShell --> Security
    Security --> Files

    style User fill:#e1f5ff,color:#1f2937
    style Intent fill:#fff4e1,color:#1f2937
    style Context fill:#fff4e1,color:#1f2937
    style LLM fill:#f3e5f5,color:#1f2937
    style Security fill:#ffebee,color:#1f2937
    style Files fill:#c8e6c9,color:#1f2937
```

## TDD Workflow (Outer + Inner Loop)

```mermaid
flowchart TD
    subgraph "Outer Loop - Acceptance Tests"
        AcceptanceRed[🔴 Write Failing<br/>Acceptance Test<br/>Ginkgo/Playwright]
        AcceptanceGreen[🟢 Acceptance Test<br/>Passes]
    end

    subgraph "Inner Loop - Unit Tests"
        UnitRed[🔴 Write Failing<br/>Unit Test<br/>Ginkgo/Gomega]
        Implementation[🟢 Write Minimal<br/>Implementation]
        UnitGreen[🟢 Unit Test<br/>Passes]
        Refactor[♻️ Refactor<br/>Improve code]
    end

    AcceptanceRed --> UnitRed
    UnitRed --> Implementation
    Implementation --> UnitGreen
    UnitGreen --> TestAcceptance{Acceptance<br/>passes?}

    TestAcceptance -->|No| UnitRed
    TestAcceptance -->|Yes| AcceptanceGreen

    AcceptanceGreen --> Refactor
    Refactor --> Done([✅ Feature Complete])

    style AcceptanceRed fill:#ffcdd2,color:#1f2937
    style UnitRed fill:#ffcdd2,color:#1f2937
    style Implementation fill:#fff9c4,color:#1f2937
    style UnitGreen fill:#c8e6c9,color:#1f2937
    style AcceptanceGreen fill:#c8e6c9,color:#1f2937
    style Refactor fill:#e1bee7,color:#1f2937
    style Done fill:#c8e6c9,color:#1f2937
```

## Deployment Models

```mermaid
graph TB
    subgraph "Single Binary Deployment"
        Binary[cortex binary<br/>50MB, all-in-one]
        Binary --> Edge1[Raspberry Pi]
        Binary --> Edge2[Developer Laptop]
        Binary --> Edge3[CI/CD Pipeline]
    end

    subgraph "Kubernetes Deployment"
        K8s[Kubernetes Cluster]
        K8s --> API[API Pods<br/>3 replicas]
        K8s --> UI[UI Pods<br/>2 replicas]
        K8s --> Worker[Worker Pods<br/>Auto-scaling]
        K8s --> DB[(PostgreSQL)]
        K8s --> Storage[(S3/MinIO)]
    end

    subgraph "Features by Deployment"
        Binary -.->|Supports| F1[Core execution<br/>Local neurons<br/>CLI interface]
        K8s -.->|Supports| F2[All features<br/>Web UI<br/>Fleet management<br/>Multi-user]
    end

    style Binary fill:#e1f5ff,color:#1f2937
    style K8s fill:#f3e5f5,color:#1f2937
    style Edge1 fill:#c8e6c9,color:#1f2937
    style Edge2 fill:#c8e6c9,color:#1f2937
    style Edge3 fill:#c8e6c9,color:#1f2937
```

## Neuron File Structure

```mermaid
graph TD
    Neuron[Neuron: check_nginx] --> Config[config.yml<br/>Metadata]
    Neuron --> Script[run.sh<br/>Executable]
    Neuron --> Readme[README.md<br/>Documentation]
    Neuron --> Tests[run_test.sh<br/>Tests]

    Config --> Name[name: check_nginx]
    Config --> Type[type: check]
    Config --> ExecFile[exec_file: run.sh]
    Config --> PreExec[pre_exec_debug: message]
    Config --> PostSuccess[post_exec_success_debug]
    Config --> PostFail[post_exec_fail_debug:<br/>  120: 'Nginx not running'<br/>  121: 'Nginx not installed']

    Script --> Shebang[#!/bin/bash]
    Script --> ErrorHandling[set -euo pipefail]
    Script --> Logic[Check logic]
    Script --> ExitCode[exit with code]

    style Neuron fill:#e1f5ff,color:#1f2937
    style Config fill:#fff4e1,color:#1f2937
    style Script fill:#c8e6c9,color:#1f2937
```

## Security Scanning Pipeline

```mermaid
flowchart TD
    subgraph "Triggers"
        Push[Git Push]
        PR[Pull Request]
        Schedule[Daily Schedule<br/>2 AM UTC]
    end

    subgraph "Go Security"
        GoVuln[govulncheck<br/>Go vulnerabilities]
        Gosec[gosec<br/>Security issues]
        GoMod[go mod verify<br/>Module integrity]
    end

    subgraph "npm Security"
        NpmAudit[npm audit<br/>npm vulnerabilities]
    end

    subgraph "Code Analysis"
        CodeQL[CodeQL<br/>Semantic analysis]
        Trivy[Trivy<br/>Container/FS scan]
    end

    subgraph "Results"
        SARIF[Upload SARIF<br/>to Security tab]
        Artifacts[Upload artifacts<br/>for review]
        Block{Block<br/>on critical?}
    end

    Push --> GoVuln
    PR --> GoVuln
    Schedule --> GoVuln

    GoVuln --> GoMod
    GoVuln --> Gosec
    GoVuln --> NpmAudit
    GoVuln --> CodeQL
    GoVuln --> Trivy

    Gosec --> SARIF
    CodeQL --> SARIF
    Trivy --> SARIF
    NpmAudit --> Artifacts

    SARIF --> Block
    Block -->|Critical| Fail[❌ Block PR/Push]
    Block -->|Lower| Pass[✅ Allow with warning]

    style Push fill:#e1f5ff,color:#1f2937
    style PR fill:#e1f5ff,color:#1f2937
    style Schedule fill:#e1f5ff,color:#1f2937
    style GoVuln fill:#c8e6c9,color:#1f2937
    style CodeQL fill:#fff4e1,color:#1f2937
    style SARIF fill:#f3e5f5,color:#1f2937
    style Fail fill:#ffcdd2,color:#1f2937
    style Pass fill:#c8e6c9,color:#1f2937
```

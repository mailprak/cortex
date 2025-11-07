# Technical Specification: AI-Powered Neuron Generation

**Version:** 1.0
**Status:** Draft
**Author:** Cortex Core Team
**Last Updated:** 2025-01-07

---

## Table of Contents

1. [Overview](#overview)
2. [Goals & Non-Goals](#goals--non-goals)
3. [Architecture](#architecture)
4. [Detailed Design](#detailed-design)
5. [API Specification](#api-specification)
6. [Implementation Plan](#implementation-plan)
7. [Testing Strategy](#testing-strategy)
8. [Security & Privacy](#security--privacy)
9. [Performance Requirements](#performance-requirements)
10. [Migration & Rollout](#migration--rollout)

---

## Overview

### Problem Statement

Current Cortex workflow requires manual creation of neurons with significant YAML and shell scripting knowledge:

**Pain Points:**
- High barrier to entry for new users
- Time-consuming manual YAML configuration
- Requires deep knowledge of exit codes and error handling
- No intelligent suggestions or best practices enforcement
- Repetitive boilerplate creation

**Current Workflow:**
```bash
# User must manually:
1. Run: cortex create-neuron check_disk_space
2. Edit neuron.yaml (understand schema)
3. Write run.sh (shell scripting knowledge required)
4. Define exit codes and error messages
5. Test and debug
Time: 15-30 minutes per neuron
```

### Proposed Solution

AI-powered neuron generation that converts natural language descriptions into production-ready neurons with best practices built-in.

**Target Workflow:**
```bash
cortex generate neuron "Check if disk space is above 80% and send Slack alert"
# → AI generates complete neuron in < 5 seconds
# → User reviews, tweaks if needed, done
Time: 30 seconds to 2 minutes
```

---

## Goals & Non-Goals

### Goals

✅ **Primary Goals:**
1. Reduce neuron creation time from 15-30 minutes to < 2 minutes
2. Lower barrier to entry for new Cortex users
3. Enforce best practices automatically (error handling, exit codes, logging)
4. Support multiple execution environments (bash, PowerShell, Python)
5. Enable context-aware generation (learn from existing neurons)

✅ **Secondary Goals:**
1. Support iterative refinement ("make it also check CPU usage")
2. Generate test cases automatically
3. Suggest related neurons from community
4. Create documentation inline

### Non-Goals

❌ **Explicitly Out of Scope:**
1. Real-time execution (AI generates, doesn't execute)
2. Automatic deployment to production
3. Complex multi-neuron synapse generation (Phase 2)
4. Training custom LLM models
5. On-device AI inference (requires API initially)

---

## Architecture

### System Context Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Cortex CLI                            │
│                                                               │
│  ┌──────────────┐    ┌─────────────────┐   ┌─────────────┐ │
│  │   User Input │───▶│  AI Generator   │──▶│   Validator │ │
│  │  (natural    │    │   Controller    │   │  & Linter   │ │
│  │   language)  │    └─────────────────┘   └─────────────┘ │
│  └──────────────┘             │                     │        │
│                                ▼                     ▼        │
│                    ┌────────────────────┐  ┌──────────────┐ │
│                    │  Template Engine   │  │  File Writer │ │
│                    │  (Go templates)    │  │              │ │
│                    └────────────────────┘  └──────────────┘ │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
            ┌────────────────────────┐
            │   LLM Provider API     │
            │  ┌──────────────────┐  │
            │  │ OpenAI GPT-4     │  │
            │  │ Anthropic Claude │  │
            │  │ Local Ollama     │  │
            │  │ Azure OpenAI     │  │
            │  └──────────────────┘  │
            └────────────────────────┘
                         │
                         ▼
            ┌────────────────────────┐
            │  Context Sources       │
            │  - Existing neurons    │
            │  - Community templates │
            │  - System metadata     │
            └────────────────────────┘
```

### Component Architecture

```go
// High-level component structure

pkg/
├── ai/
│   ├── generator.go          // Main AI generation orchestrator
│   ├── providers/
│   │   ├── openai.go         // OpenAI GPT integration
│   │   ├── anthropic.go      // Claude integration
│   │   ├── ollama.go         // Local Ollama integration
│   │   └── interface.go      // LLM provider interface
│   ├── context/
│   │   ├── collector.go      // Gather context from system
│   │   └── embeddings.go     // Vector search for similar neurons
│   └── prompts/
│       ├── templates.go      // System prompts
│       └── examples.go       // Few-shot examples
├── neuron/
│   ├── validator.go          // Validate generated neurons
│   ├── linter.go            // Best practice checks
│   └── tester.go            // Auto-test generation
├── templates/
│   └── neuron_templates/    // Go templates for generation
└── config/
    └── ai_config.go         // AI provider configuration
```

---

## Detailed Design

### 1. User Interface Design

#### CLI Commands

```bash
# Basic generation
cortex generate neuron "description"
cortex gen neuron "description"  # alias

# With options
cortex generate neuron "check disk space" \
  --language bash \
  --platform linux \
  --context ./existing-neurons/ \
  --provider openai \
  --model gpt-4 \
  --output ./neurons/check_disk

# Interactive mode
cortex generate neuron --interactive

# From template with AI enhancement
cortex generate neuron --template monitoring/disk \
  --enhance "also check inode usage"

# Dry run (preview only)
cortex generate neuron "description" --dry-run
```

#### Interactive Prompts

```bash
$ cortex generate neuron --interactive

🤖 Cortex AI Neuron Generator

What should this neuron do?
> Check if PostgreSQL replication lag is under 5 seconds

Which execution environment?
[1] Bash (Linux/macOS)
[2] PowerShell (Windows)
[3] Python 3
[4] Auto-detect
> 1

Should this neuron mutate anything? (check vs mutate)
[1] Check only (read-only operations)
[2] Mutate (modify system state)
> 1

Generate tests automatically? [Y/n] y

Analyzing requirements...
✓ Detected: PostgreSQL, monitoring, replication
✓ Found 3 similar neurons in community
✓ Generating code...

Generated neuron: check_postgres_replication_lag

📁 Files created:
  - neurons/check_postgres_replication_lag/neuron.yaml
  - neurons/check_postgres_replication_lag/run.sh
  - neurons/check_postgres_replication_lag/README.md
  - neurons/check_postgres_replication_lag/run_test.sh

📝 Summary:
  - Exit Code 0: Replication lag OK (< 5s)
  - Exit Code 120: Replication lag warning (5-10s)
  - Exit Code 121: Replication lag critical (> 10s)
  - Exit Code 122: PostgreSQL not reachable

Review and test? [Y/n] y
```

### 2. AI Generation Pipeline

#### Pipeline Stages

```
┌─────────────────────────────────────────────────────────────────┐
│                     STAGE 1: INTENT ANALYSIS                     │
│  Input: Natural language description                             │
│  Output: Structured intent object                                │
│  Process:                                                         │
│    - Extract: What to check/do                                   │
│    - Extract: Success/failure conditions                         │
│    - Extract: Technologies involved                              │
│    - Classify: Check vs Mutate                                   │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                    STAGE 2: CONTEXT GATHERING                    │
│  Input: Intent object                                            │
│  Output: Enriched context                                        │
│  Process:                                                         │
│    - Search similar neurons (vector similarity)                  │
│    - Load relevant templates                                     │
│    - Gather system info (OS, tools available)                    │
│    - Check community best practices                              │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                     STAGE 3: CODE GENERATION                     │
│  Input: Intent + Context                                         │
│  Output: Raw generated code                                      │
│  Process:                                                         │
│    - Construct system prompt                                     │
│    - Add few-shot examples                                       │
│    - Call LLM API                                                │
│    - Parse structured response                                   │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                   STAGE 4: VALIDATION & LINTING                  │
│  Input: Generated code                                           │
│  Output: Validated, linted code                                  │
│  Process:                                                         │
│    - YAML schema validation                                      │
│    - Shell script syntax check (shellcheck)                      │
│    - Security scan (no hardcoded secrets)                        │
│    - Best practice enforcement                                   │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                      STAGE 5: TEST GENERATION                    │
│  Input: Validated neuron                                         │
│  Output: Test suite                                              │
│  Process:                                                         │
│    - Generate happy path test                                    │
│    - Generate error condition tests                              │
│    - Create mock data if needed                                  │
└──────────────────────────────┬──────────────────────────────────┘
                               │
┌──────────────────────────────▼──────────────────────────────────┐
│                       STAGE 6: FILE WRITING                      │
│  Input: Complete neuron package                                  │
│  Output: Files on disk                                           │
│  Process:                                                         │
│    - Create directory structure                                  │
│    - Write neuron.yaml                                           │
│    - Write execution script (run.sh/run.ps1)                     │
│    - Write README.md                                             │
│    - Write tests (if enabled)                                    │
└──────────────────────────────────────────────────────────────────┘
```

### 3. Data Structures

#### Intent Object

```go
// pkg/ai/intent.go

type NeuronIntent struct {
    // Core description
    Description string `json:"description"`

    // Extracted metadata
    Type            NeuronType   `json:"type"`             // Check or Mutate
    Technologies    []string     `json:"technologies"`     // ["postgresql", "sql"]
    Platform        Platform     `json:"platform"`         // Linux, Windows, etc.
    Language        Language     `json:"language"`         // Bash, PowerShell, Python

    // Success/failure criteria
    SuccessCondition string      `json:"success_condition"`
    FailureScenarios []Scenario  `json:"failure_scenarios"`

    // Execution details
    ExpectedDuration time.Duration `json:"expected_duration"`
    RequiresElevated bool          `json:"requires_elevated"` // sudo/admin
    ExternalDeps     []string      `json:"external_deps"`     // psql, curl, etc.

    // User preferences
    GenerateTests    bool         `json:"generate_tests"`
    Verbose          bool         `json:"verbose"`
}

type NeuronType string
const (
    NeuronTypeCheck  NeuronType = "check"
    NeuronTypeMutate NeuronType = "mutate"
)

type Scenario struct {
    ExitCode    int    `json:"exit_code"`
    Description string `json:"description"`
    Severity    string `json:"severity"` // warning, error, critical
}
```

#### LLM Provider Interface

```go
// pkg/ai/providers/interface.go

type LLMProvider interface {
    // Generate neuron code from structured request
    Generate(ctx context.Context, req GenerateRequest) (GenerateResponse, error)

    // Health check
    Ping(ctx context.Context) error

    // Estimate cost before generation
    EstimateCost(req GenerateRequest) (float64, error)

    // Provider metadata
    Name() string
    Models() []string
}

type GenerateRequest struct {
    SystemPrompt string
    UserPrompt   string
    Temperature  float32
    MaxTokens    int
    Model        string
    Context      []Message // Few-shot examples
}

type GenerateResponse struct {
    // Structured output
    NeuronYAML   string
    ExecutionScript string
    README       string
    Tests        string

    // Metadata
    Model        string
    TokensUsed   int
    Cost         float64
    Duration     time.Duration
}

type Message struct {
    Role    string `json:"role"`    // system, user, assistant
    Content string `json:"content"`
}
```

### 4. Prompt Engineering

#### System Prompt Template

```go
// pkg/ai/prompts/templates.go

const SystemPromptTemplate = `You are an expert Site Reliability Engineer and Cortex neuron generator.

CORTEX NEURON STRUCTURE:
A neuron is a discrete, reusable infrastructure debugging task consisting of:
1. neuron.yaml - Configuration with exit codes
2. run.sh/run.ps1 - Executable script
3. README.md - Documentation

NAMING CONVENTIONS:
- Prefix "check_" for read-only operations
- Prefix "mutate_" for state-changing operations
- Use snake_case: check_disk_space, mutate_restart_service

EXIT CODE STANDARDS:
- 0: Success
- 1-99: Reserved for system errors
- 100-199: Application-specific success variants
- 200-255: Application-specific failures

BEST PRACTICES:
1. Always check prerequisites (command exists, permissions)
2. Provide clear error messages in post_exec_fail_debug
3. Use set -euo pipefail in bash scripts
4. Include timeout protection
5. Log to stderr for errors, stdout for output
6. Never hardcode secrets (use env vars)
7. Make scripts idempotent where possible

OUTPUT FORMAT:
Respond with valid JSON containing these fields:
{
  "neuron_yaml": "...",
  "execution_script": "...",
  "readme": "...",
  "exit_codes": [
    {"code": 0, "meaning": "..."},
    {"code": 120, "meaning": "..."}
  ]
}
`

const UserPromptTemplate = `Generate a Cortex neuron with the following requirements:

DESCRIPTION: {{.Description}}

TYPE: {{.Type}}
PLATFORM: {{.Platform}}
LANGUAGE: {{.Language}}

SUCCESS CONDITION: {{.SuccessCondition}}

FAILURE SCENARIOS:
{{range .FailureScenarios}}
- {{.Description}} (severity: {{.Severity}})
{{end}}

EXTERNAL DEPENDENCIES: {{.ExternalDeps}}

CONTEXT FROM SIMILAR NEURONS:
{{range .SimilarNeurons}}
Name: {{.Name}}
Exit Codes: {{.ExitCodes}}
{{end}}

Generate a production-ready neuron following all best practices.
`
```

#### Few-Shot Examples

```go
// pkg/ai/prompts/examples.go

var FewShotExamples = []Message{
    {
        Role: "user",
        Content: "Generate a neuron to check if nginx is running",
    },
    {
        Role: "assistant",
        Content: `{
  "neuron_yaml": "---\nname: check_nginx_running\ntype: check\ndescription: \"Check if nginx process is running\"\nexec_file: run.sh\npre_exec_debug: \"Checking nginx process status...\"\nassertExitStatus: [0]\npost_exec_success_debug: \"nginx is running\"\npost_exec_fail_debug:\n  120: \"nginx process not found\"\n  121: \"nginx not responding to signals\"\n",
  "execution_script": "#!/bin/bash\nset -euo pipefail\n\n# Check if nginx command exists\nif ! command -v nginx &> /dev/null; then\n    echo \"nginx command not found\" >&2\n    exit 121\nfi\n\n# Check if nginx is running\nif pgrep -x nginx > /dev/null; then\n    echo \"nginx is running (PID: $(pgrep -x nginx | head -1))\"\n    exit 0\nelse\n    echo \"nginx is not running\" >&2\n    exit 120\nfi\n",
  "readme": "# check_nginx_running\n\nChecks if nginx web server is running.\n\n## Exit Codes\n- 0: nginx is running\n- 120: nginx process not found\n- 121: nginx command not available\n\n## Prerequisites\n- nginx must be installed\n- Standard process monitoring permissions\n",
  "exit_codes": [
    {"code": 0, "meaning": "nginx is running"},
    {"code": 120, "meaning": "nginx process not found"},
    {"code": 121, "meaning": "nginx not available"}
  ]
}`,
    },
    // Add 2-3 more examples for different scenarios
}
```

### 5. Context Collection

#### Similar Neuron Search

```go
// pkg/ai/context/collector.go

type ContextCollector struct {
    neuronDir     string
    embeddingAPI  EmbeddingService
    vectorDB      VectorStore
}

func (c *ContextCollector) FindSimilarNeurons(intent NeuronIntent, limit int) ([]SimilarNeuron, error) {
    // 1. Generate embedding for user's description
    embedding, err := c.embeddingAPI.Embed(intent.Description)
    if err != nil {
        return nil, err
    }

    // 2. Search vector database
    results, err := c.vectorDB.Search(embedding, limit)
    if err != nil {
        return nil, err
    }

    // 3. Load full neuron metadata
    var similar []SimilarNeuron
    for _, result := range results {
        neuron, err := c.loadNeuron(result.ID)
        if err != nil {
            continue // Skip if can't load
        }
        similar = append(similar, SimilarNeuron{
            Name:       neuron.Name,
            ExitCodes:  neuron.ExitCodes,
            Similarity: result.Score,
        })
    }

    return similar, nil
}

// Simple implementation without vector DB (Phase 1)
func (c *ContextCollector) FindSimilarNeuronsSimple(intent NeuronIntent) ([]SimilarNeuron, error) {
    var similar []SimilarNeuron

    // Load all existing neurons
    neurons, err := c.loadAllNeurons()
    if err != nil {
        return nil, err
    }

    // Simple keyword matching
    keywords := extractKeywords(intent.Description)

    for _, neuron := range neurons {
        score := calculateSimilarity(keywords, neuron)
        if score > 0.3 { // Threshold
            similar = append(similar, SimilarNeuron{
                Name:       neuron.Name,
                ExitCodes:  neuron.ExitCodes,
                Similarity: score,
            })
        }
    }

    // Sort by similarity descending
    sort.Slice(similar, func(i, j int) bool {
        return similar[i].Similarity > similar[j].Similarity
    })

    return similar[:min(5, len(similar))], nil
}
```

---

## API Specification

### Command-Line Interface

```bash
cortex generate neuron [OPTIONS] DESCRIPTION

OPTIONS:
  -l, --language STRING      Execution language (bash|powershell|python) [default: auto]
  -p, --platform STRING      Target platform (linux|windows|darwin) [default: current]
  -t, --type STRING         Neuron type (check|mutate) [default: auto]
  -o, --output PATH         Output directory [default: ./neurons/<generated-name>]

  --provider STRING         LLM provider (openai|anthropic|ollama|azure) [default: from config]
  --model STRING            Specific model to use [default: provider default]
  --temperature FLOAT       LLM temperature (0.0-1.0) [default: 0.2]

  --context PATH            Include context from existing neurons directory
  --no-tests               Don't generate test files
  --no-readme              Don't generate README

  --interactive            Interactive mode with prompts
  --dry-run                Preview generation without writing files
  --verbose                Show detailed generation process

  --config PATH            Custom AI config file [default: ~/.cortex/ai.yaml]

EXAMPLES:
  # Basic usage
  cortex generate neuron "check if disk usage is above 80%"

  # With specific language
  cortex gen neuron "restart apache service" --language bash --type mutate

  # Using local Ollama
  cortex gen neuron "check PostgreSQL replication" --provider ollama --model llama2

  # Preview before creating
  cortex gen neuron "check API health" --dry-run
```

### Configuration File

```yaml
# ~/.cortex/ai.yaml

ai:
  # Default provider
  default_provider: openai

  # Provider configurations
  providers:
    openai:
      api_key: ${OPENAI_API_KEY}  # Use env var
      model: gpt-4-turbo
      organization: org-xxx
      timeout: 30s

    anthropic:
      api_key: ${ANTHROPIC_API_KEY}
      model: claude-3-5-sonnet-20250107
      timeout: 30s

    ollama:
      endpoint: http://localhost:11434
      model: llama2
      timeout: 60s

    azure:
      api_key: ${AZURE_OPENAI_KEY}
      endpoint: https://your-resource.openai.azure.com
      deployment: gpt-4
      api_version: 2024-02-15-preview

  # Generation settings
  generation:
    temperature: 0.2          # Conservative for code generation
    max_tokens: 2048
    include_similar_neurons: true
    max_similar: 3

  # Feature flags
  features:
    generate_tests: true
    generate_readme: true
    security_scan: true
    auto_lint: true

  # Validation
  validation:
    shellcheck_enabled: true
    yaml_schema_check: true

  # Privacy
  privacy:
    send_telemetry: false     # Don't send usage data
    anonymize_paths: true     # Remove sensitive paths from prompts
```

### Go API

```go
// pkg/ai/generator.go

package ai

import (
    "context"
    "time"
)

// Generator orchestrates AI-powered neuron generation
type Generator struct {
    provider     LLMProvider
    validator    *neuron.Validator
    contextCol   *context.Collector
    templateEng  *templates.Engine
}

// NewGenerator creates a new AI generator
func NewGenerator(config Config) (*Generator, error) {
    provider, err := providers.NewProvider(config.Provider)
    if err != nil {
        return nil, err
    }

    return &Generator{
        provider:    provider,
        validator:   neuron.NewValidator(),
        contextCol:  context.NewCollector(config.NeuronDir),
        templateEng: templates.NewEngine(),
    }, nil
}

// Generate creates a neuron from natural language description
func (g *Generator) Generate(ctx context.Context, req GenerateRequest) (*GenerateResult, error) {
    // Stage 1: Analyze intent
    intent, err := g.analyzeIntent(req.Description)
    if err != nil {
        return nil, fmt.Errorf("intent analysis failed: %w", err)
    }

    // Stage 2: Gather context
    contextData, err := g.contextCol.Gather(intent)
    if err != nil {
        return nil, fmt.Errorf("context gathering failed: %w", err)
    }

    // Stage 3: Generate code
    generated, err := g.generateCode(ctx, intent, contextData)
    if err != nil {
        return nil, fmt.Errorf("code generation failed: %w", err)
    }

    // Stage 4: Validate
    if err := g.validator.Validate(generated); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // Stage 5: Generate tests (if enabled)
    if req.GenerateTests {
        generated.Tests, err = g.generateTests(generated)
        if err != nil {
            return nil, fmt.Errorf("test generation failed: %w", err)
        }
    }

    return &GenerateResult{
        Intent:    intent,
        Generated: generated,
        Metadata:  g.buildMetadata(generated),
    }, nil
}

// GenerateRequest contains all parameters for generation
type GenerateRequest struct {
    Description   string
    Language      Language
    Platform      Platform
    Type          NeuronType
    GenerateTests bool
    ContextPath   string
    DryRun        bool
}

// GenerateResult contains the complete generated neuron
type GenerateResult struct {
    Intent    NeuronIntent
    Generated GeneratedNeuron
    Metadata  GenerationMetadata
}

type GeneratedNeuron struct {
    NeuronYAML      string
    ExecutionScript string
    README          string
    Tests           string
    ExitCodes       []ExitCodeMapping
}

type GenerationMetadata struct {
    Model          string
    TokensUsed     int
    Cost           float64
    Duration       time.Duration
    SimilarNeurons []string
}
```

---

## Implementation Plan

### Phase 1: MVP (Weeks 1-4)

**Week 1: Infrastructure**
- [ ] Set up project structure
- [ ] Implement LLM provider interface
- [ ] Add OpenAI provider implementation
- [ ] Create configuration system
- [ ] Write unit tests for providers

**Week 2: Core Generation**
- [ ] Implement intent analysis (basic)
- [ ] Build prompt templates
- [ ] Create code generation pipeline
- [ ] Add YAML/shell validation
- [ ] Write integration tests

**Week 3: CLI & UX**
- [ ] Build CLI commands
- [ ] Add interactive mode
- [ ] Implement dry-run preview
- [ ] Create progress indicators
- [ ] Add error handling

**Week 4: Testing & Documentation**
- [ ] End-to-end testing
- [ ] User documentation
- [ ] Example gallery
- [ ] Performance testing
- [ ] Beta release

### Phase 2: Enhancements (Weeks 5-8)

**Week 5-6: Additional Providers**
- [ ] Anthropic Claude integration
- [ ] Ollama local provider
- [ ] Azure OpenAI support
- [ ] Provider auto-fallback

**Week 7-8: Intelligence**
- [ ] Context from existing neurons
- [ ] Vector similarity search
- [ ] Community template integration
- [ ] Auto-test generation

### Phase 3: Advanced Features (Weeks 9-12)

**Week 9-10: Refinement**
- [ ] Iterative enhancement
- [ ] Multi-neuron generation
- [ ] Synapse suggestion
- [ ] Performance optimization

**Week 11-12: Production Ready**
- [ ] Security hardening
- [ ] Telemetry (optional)
- [ ] Cost tracking
- [ ] 1.0 Release

---

## Testing Strategy

### Unit Tests

```go
// pkg/ai/generator_test.go

func TestGenerator_AnalyzeIntent(t *testing.T) {
    tests := []struct {
        name        string
        description string
        want        NeuronIntent
        wantErr     bool
    }{
        {
            name:        "disk space check",
            description: "check if disk usage is above 80%",
            want: NeuronIntent{
                Type:             NeuronTypeCheck,
                Technologies:     []string{"disk", "filesystem"},
                SuccessCondition: "disk usage < 80%",
            },
            wantErr: false,
        },
        {
            name:        "service restart",
            description: "restart nginx service",
            want: NeuronIntent{
                Type:         NeuronTypeMutate,
                Technologies: []string{"nginx", "systemd"},
            },
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            g := NewGenerator(testConfig)
            got, err := g.analyzeIntent(tt.description)

            if (err != nil) != tt.wantErr {
                t.Errorf("analyzeIntent() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if got.Type != tt.want.Type {
                t.Errorf("analyzeIntent() Type = %v, want %v", got.Type, tt.want.Type)
            }
        })
    }
}
```

### Integration Tests

```go
func TestGenerator_Generate_E2E(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Requires valid API key
    apiKey := os.Getenv("OPENAI_API_KEY")
    if apiKey == "" {
        t.Skip("OPENAI_API_KEY not set")
    }

    ctx := context.Background()
    g := setupTestGenerator(t, apiKey)

    result, err := g.Generate(ctx, GenerateRequest{
        Description:   "check if nginx is running",
        GenerateTests: true,
    })

    require.NoError(t, err)
    require.NotEmpty(t, result.Generated.NeuronYAML)
    require.NotEmpty(t, result.Generated.ExecutionScript)

    // Validate YAML parses correctly
    var neuronConfig neuron.Config
    err = yaml.Unmarshal([]byte(result.Generated.NeuronYAML), &neuronConfig)
    require.NoError(t, err)

    // Validate script syntax
    err = validateShellScript(result.Generated.ExecutionScript)
    require.NoError(t, err)
}
```

### Acceptance Tests

```bash
#!/bin/bash
# test/acceptance/test_ai_generation.sh

set -euo pipefail

echo "=== Acceptance Test: AI Neuron Generation ==="

# Test 1: Basic generation
echo "Test 1: Generate simple check neuron"
cortex generate neuron "check if port 8080 is open" \
    --output /tmp/test-neuron \
    --dry-run > /tmp/gen-output.txt

grep -q "neuron.yaml" /tmp/gen-output.txt || exit 1
echo "✓ Test 1 passed"

# Test 2: Interactive mode
echo "Test 2: Interactive generation"
echo -e "check disk space\n1\n1\ny\nn\n" | cortex generate neuron --interactive
echo "✓ Test 2 passed"

# Test 3: Actual file creation
echo "Test 3: Create real neuron"
cortex generate neuron "check nginx status" --output /tmp/nginx-check
test -f /tmp/nginx-check/neuron.yaml || exit 1
test -f /tmp/nginx-check/run.sh || exit 1
test -x /tmp/nginx-check/run.sh || exit 1
echo "✓ Test 3 passed"

# Test 4: Validate generated neuron works
echo "Test 4: Execute generated neuron"
cd /tmp/nginx-check
chmod +x run.sh
./run.sh && echo "✓ Test 4 passed" || echo "⚠ Test 4: neuron execution failed (expected if nginx not running)"

echo "=== All acceptance tests passed ==="
```

---

## Security & Privacy

### Security Considerations

#### 1. API Key Protection

```go
// Never log API keys
func (p *OpenAIProvider) Generate(ctx context.Context, req GenerateRequest) (*GenerateResponse, error) {
    // ✅ CORRECT: Use context for API key
    apiKey := ctx.Value("api_key").(string)

    // ❌ WRONG: Never log the full key
    log.Printf("Using API key: %s", maskAPIKey(apiKey)) // Only log masked version

    // Make request...
}

func maskAPIKey(key string) string {
    if len(key) < 8 {
        return "****"
    }
    return key[:4] + "****" + key[len(key)-4:]
}
```

#### 2. Prompt Injection Prevention

```go
// Sanitize user input to prevent prompt injection
func sanitizeDescription(desc string) string {
    // Remove potential prompt injection attempts
    dangerous := []string{
        "ignore previous instructions",
        "system prompt",
        "you are now",
    }

    cleaned := strings.ToLower(desc)
    for _, d := range dangerous {
        if strings.Contains(cleaned, d) {
            return "" // Reject suspicious input
        }
    }

    // Limit length
    if len(desc) > 500 {
        desc = desc[:500]
    }

    return desc
}
```

#### 3. Generated Code Validation

```go
// Scan generated code for security issues
func (v *Validator) SecurityScan(script string) error {
    // Check for hardcoded credentials
    credentialPatterns := []string{
        `password\s*=\s*["'][^"']+["']`,
        `api_key\s*=\s*["'][^"']+["']`,
        `secret\s*=\s*["'][^"']+["']`,
    }

    for _, pattern := range credentialPatterns {
        if matched, _ := regexp.MatchString(pattern, script); matched {
            return fmt.Errorf("security: hardcoded credential detected")
        }
    }

    // Check for dangerous commands
    dangerousCommands := []string{
        "rm -rf /",
        ":(){ :|:& };:",  // Fork bomb
        "curl | bash",
    }

    for _, cmd := range dangerousCommands {
        if strings.Contains(script, cmd) {
            return fmt.Errorf("security: dangerous command detected: %s", cmd)
        }
    }

    return nil
}
```

### Privacy Considerations

#### 1. Data Minimization

```go
// Remove sensitive data before sending to LLM
func (g *Generator) anonymizeContext(ctx *ContextData) *ContextData {
    anonymized := *ctx

    // Remove absolute paths
    anonymized.NeuronDir = "/path/to/neurons"

    // Remove hostnames
    anonymized.Hostname = "localhost"

    // Remove IPs
    re := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
    anonymized.Description = re.ReplaceAllString(anonymized.Description, "X.X.X.X")

    return &anonymized
}
```

#### 2. Opt-in Telemetry

```yaml
# Default: telemetry disabled
ai:
  privacy:
    send_telemetry: false  # User must explicitly enable
    anonymize_paths: true

telemetry:
  enabled: false
  # If enabled, only send:
  endpoint: https://telemetry.cortex.dev/v1/events
  send_data:
    - generation_success_count
    - generation_failure_count
    - avg_generation_time
  # Never send:
    # - User descriptions (PII)
    # - Generated code
    # - API keys
    # - Paths
```

---

## Performance Requirements

### Latency Requirements

| Operation | Target | Maximum | Notes |
|-----------|--------|---------|-------|
| Simple neuron generation | < 5s | 10s | "check disk space" |
| Complex neuron generation | < 10s | 20s | Multi-step, specific tech |
| Context gathering | < 1s | 3s | From local neurons |
| Validation | < 100ms | 500ms | YAML + shell syntax |
| Total end-to-end | < 6s | 15s | From command to files |

### Resource Requirements

```yaml
memory:
  baseline: 50MB      # Binary + minimal runtime
  per_generation: 5MB # Additional per concurrent generation
  maximum: 200MB      # Total cap

cpu:
  generation: 1 core   # Can parallelize multiple generations
  validation: 0.1 core # Lightweight

disk:
  cache: 100MB        # Template cache, similar neuron index
  per_neuron: 10KB    # Each generated neuron

network:
  bandwidth: 10KB/s   # LLM API requests (compressed)
  concurrent: 5       # Max parallel API calls
```

### Scalability Targets

- **Sequential**: 10 neurons/minute (user typing descriptions)
- **Batch**: 100 neurons/hour (from template list)
- **Concurrent users**: 50 simultaneous generations

---

## Migration & Rollout

### Rollout Plan

#### Phase 1: Alpha (Week 4)
```bash
# Feature flag controlled
cortex generate neuron --enable-ai-alpha "description"

# Limited to:
- Internal team testing (10 users)
- OpenAI provider only
- Manual approval before file write
```

#### Phase 2: Beta (Week 6)
```bash
# Public beta
cortex generate neuron "description"  # Works by default

# Available to:
- Public GitHub repository
- Opt-in beta testers (100 users)
- Multiple providers (OpenAI, Anthropic)
- Auto-write with --yes flag
```

#### Phase 3: GA (Week 12)
```bash
# General availability
cortex generate neuron "description"

# Available to:
- All users
- Full documentation
- Community templates
- Production-ready
```

### Backward Compatibility

```bash
# Existing workflow still works
cortex create-neuron my_neuron  # Manual creation unchanged

# New AI workflow is additive
cortex generate neuron "description"  # New capability

# Users can mix both approaches
cortex create-neuron template_neuron
cortex generate neuron "enhance template with AI" --template template_neuron
```

### Migration Guide

```markdown
# Migrating from Manual to AI-Assisted Neuron Creation

## Before (Manual)
1. Run `cortex create-neuron check_disk`
2. Edit `check_disk/neuron.yaml` manually
3. Write `check_disk/run.sh` from scratch
4. Test and debug
5. Time: 15-30 minutes

## After (AI-Assisted)
1. Run `cortex generate neuron "check disk usage above 80%"`
2. Review generated files
3. Tweak if needed
4. Time: 1-2 minutes

## Hybrid Approach
- Use AI for initial scaffolding
- Manually refine complex logic
- Best of both worlds
```

---

## Appendix

### A. Example Generated Neuron

**Input:**
```bash
cortex generate neuron "Check if PostgreSQL replication lag is under 5 seconds"
```

**Output - neuron.yaml:**
```yaml
---
name: check_postgres_replication_lag
type: check
description: "Check if PostgreSQL replication lag is under 5 seconds"
exec_file: run.sh
pre_exec_debug: "Checking PostgreSQL replication lag..."
assertExitStatus: [0]
post_exec_success_debug: "PostgreSQL replication lag is healthy (< 5 seconds)"
post_exec_fail_debug:
  120: "Replication lag is in warning range (5-10 seconds)"
  121: "Replication lag is critical (> 10 seconds)"
  122: "Unable to connect to PostgreSQL"
  123: "Replication is not configured"
```

**Output - run.sh:**
```bash
#!/bin/bash
set -euo pipefail

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo "Error: psql command not found. Please install postgresql-client." >&2
    exit 122
fi

# Database connection parameters (from environment)
PGHOST=${PGHOST:-localhost}
PGPORT=${PGPORT:-5432}
PGUSER=${PGUSER:-postgres}
PGDATABASE=${PGDATABASE:-postgres}

# Query replication lag
LAG_QUERY="SELECT EXTRACT(EPOCH FROM (now() - pg_last_xact_replay_timestamp())) AS lag_seconds;"

# Execute query
LAG=$(psql -h "$PGHOST" -p "$PGPORT" -U "$PGUSER" -d "$PGDATABASE" \
    -t -c "$LAG_QUERY" 2>/dev/null) || {
    echo "Error: Failed to query PostgreSQL. Check connection parameters." >&2
    exit 122
}

# Check if replication is configured (null = not a replica)
if [ -z "$LAG" ] || [ "$LAG" == "" ]; then
    echo "Warning: This server is not configured as a PostgreSQL replica." >&2
    exit 123
fi

# Convert to integer for comparison
LAG_INT=$(echo "$LAG" | awk '{print int($1)}')

# Evaluate lag thresholds
if [ "$LAG_INT" -lt 5 ]; then
    echo "PostgreSQL replication lag: ${LAG_INT}s (healthy)"
    exit 0
elif [ "$LAG_INT" -lt 10 ]; then
    echo "PostgreSQL replication lag: ${LAG_INT}s (warning)" >&2
    exit 120
else
    echo "PostgreSQL replication lag: ${LAG_INT}s (critical)" >&2
    exit 121
fi
```

**Output - README.md:**
```markdown
# check_postgres_replication_lag

Checks if PostgreSQL replication lag is under 5 seconds.

## Description

This neuron monitors the replication lag on a PostgreSQL replica server by querying
`pg_last_xact_replay_timestamp()` and comparing it to the current time.

## Prerequisites

- `psql` (PostgreSQL client) must be installed
- PostgreSQL server must be configured as a replica
- User must have permission to query replication status

## Environment Variables

- `PGHOST`: PostgreSQL host (default: localhost)
- `PGPORT`: PostgreSQL port (default: 5432)
- `PGUSER`: PostgreSQL user (default: postgres)
- `PGDATABASE`: Database name (default: postgres)
- `PGPASSWORD`: PostgreSQL password (set if required)

## Exit Codes

- **0**: Replication lag is healthy (< 5 seconds)
- **120**: Replication lag is in warning range (5-10 seconds)
- **121**: Replication lag is critical (> 10 seconds)
- **122**: Unable to connect to PostgreSQL
- **123**: Server is not configured as a replica

## Usage

```bash
# Basic usage
cortex exec check_postgres_replication_lag

# With custom connection
export PGHOST=replica.example.com
export PGUSER=monitoring
export PGPASSWORD=secret
cortex exec check_postgres_replication_lag
```

## Integration with Synapse

```yaml
# Example synapse for database health check
definition:
  - neuron: check_postgres_replication_lag
    config:
      path: /path/to/neurons/check_postgres_replication_lag
      fix:
        120: alert_on_call_team
        121: initiate_failover_procedure
```

## Monitoring Best Practices

- Run this check every 1-5 minutes
- Alert on exit code 120 (warning)
- Page on-call for exit code 121 (critical)
- Investigate exit code 123 (misconfiguration)
```

### B. Cost Estimation

**OpenAI GPT-4 Turbo Pricing (as of 2025):**
- Input: $0.01 per 1K tokens
- Output: $0.03 per 1K tokens

**Typical Generation:**
```
System prompt: 500 tokens
User prompt: 300 tokens
Context: 200 tokens
Total input: 1,000 tokens = $0.01

Generated output: 800 tokens = $0.024

Total per neuron: ~$0.034 (3.4 cents)
```

**Cost at Scale:**
- 10 neurons/day: $0.34/day = $10/month
- 100 neurons/day: $3.40/day = $102/month
- 1000 neurons/day: $34/day = $1,020/month

**Cost Optimization:**
- Use Ollama (local, free) for development
- Cache similar neuron context
- Batch multiple requests
- Use GPT-3.5 for simple neurons (90% cheaper)

### C. Future Enhancements

**Post-1.0 Features:**
1. Multi-neuron synapse generation
2. Visual neuron builder (web UI)
3. Fine-tuned model on Cortex neurons
4. Auto-optimization of existing neurons
5. Natural language querying of neuron history
6. Integration with monitoring systems
7. Automated documentation generation
8. Cross-platform neuron translation (bash ↔ PowerShell)

---

**Document Status:** Ready for Review
**Next Steps:** Architecture review, prototype implementation
**Estimated Effort:** 12 weeks (1 engineer)
**Target Release:** Q2 2025

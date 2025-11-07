# Technical Specification: Cortex Web UI & Kubernetes Deployment

**Version:** 1.0
**Status:** Draft
**Author:** Cortex Core Team
**Last Updated:** 2025-01-07

---

## Table of Contents

1. [Overview](#overview)
2. [Design Philosophy](#design-philosophy)
3. [User Research & Personas](#user-research--personas)
4. [Architecture](#architecture)
5. [Feature Specifications](#feature-specifications)
6. [UI/UX Design Patterns](#uiux-design-patterns)
7. [Kubernetes Deployment](#kubernetes-deployment)
8. [Technical Stack](#technical-stack)
9. [Implementation Roadmap](#implementation-roadmap)
10. [Accessibility & Performance](#accessibility--performance)

---

## Overview

### Problem Statement

While Cortex's CLI is powerful and lightweight, many users need:
- **Visual workflow building** - Drag-and-drop neuron orchestration
- **Real-time monitoring** - Live execution status and logs
- **Team collaboration** - Shared dashboards and runbook libraries
- **Quick troubleshooting** - Historical execution analysis
- **Easier onboarding** - Visual guides for CLI concepts

**Current Gap:**
- CLI-only interface has steep learning curve for non-terminal users
- No way to visualize neuron execution flow
- Difficult to share debugging knowledge across teams
- No centralized view of fleet-wide neuron execution (edge devices)

### Solution: Cortex Web UI

A **lightweight, optional web interface** that complements the CLI without replacing it:

**Core Principles:**
1. **CLI First** - UI is a convenience layer, not required
2. **Lightweight** - < 10MB container, runs alongside binary
3. **Real-time** - WebSocket-based live updates
4. **Offline Capable** - Progressive Web App (PWA)
5. **Mobile Responsive** - Works on tablets/phones for on-call scenarios

**Deployment Modes:**
- **Local Binary**: `cortex ui --port 8080` (single-user development)
- **Kubernetes Service**: Multi-user production deployment with auth
- **Docker Compose**: Team deployments with persistent storage

---

## Design Philosophy

### Inspired By Best-in-Class Tools

**From Grafana:**
- ✅ Dashboard-first navigation
- ✅ Time-range selectors for historical analysis
- ✅ Panel-based layouts with drill-down
- ✅ Dark/light theme toggle

**From Portainer:**
- ✅ Clear resource hierarchy (clusters → nodes → neurons)
- ✅ Quick actions on hover
- ✅ Templates/presets for common tasks
- ✅ Visual status indicators

**From Kubernetes Dashboard:**
- ✅ Real-time resource updates
- ✅ Log streaming with search/filter
- ✅ YAML editor with validation
- ✅ Resource creation wizards

**From GitHub Actions:**
- ✅ Workflow visualization (DAG view)
- ✅ Live execution logs
- ✅ Re-run failed jobs
- ✅ Artifacts/output downloads

### Unique Cortex UI Principles

1. **AI-Assisted Everything**
   - Natural language search: "Show failed disk checks last 24h"
   - AI suggestions: "This neuron often fails with exit 120, try X"
   - Auto-remediation: "Fix detected, apply automatically?"

2. **Edge-Aware**
   - Fleet view: See all deployed Cortex instances (IoT devices)
   - Offline indicator: Which edge nodes are unreachable
   - Bandwidth-conscious: Minimal data transfer for remote access

3. **Shell Heritage**
   - Embedded terminal for power users
   - Copy commands directly from UI
   - Export workflows as CLI commands

4. **Progressive Disclosure**
   - Start simple (run a neuron), reveal complexity gradually
   - Beginner mode vs Expert mode toggle
   - Contextual help at every step

---

## User Research & Personas

### Persona 1: "Sarah - Senior SRE"

**Profile:**
- 8 years experience, manages 500+ servers
- Prefers CLI but needs team collaboration
- On-call 1 week/month, responds via phone

**Needs:**
- Quick access to historical neuron runs
- Mobile-friendly dashboard for on-call
- Share debugging runbooks with team
- Audit trail of who ran what when

**UI Requirements:**
- Mobile-responsive design
- Role-based access control (RBAC)
- Exportable reports
- Slack/PagerDuty integrations

### Persona 2: "Alex - DevOps Engineer (Junior)"

**Profile:**
- 2 years experience, learning infrastructure
- Comfortable with GUIs, less with CLIs
- Needs visual aids to understand concepts

**Needs:**
- Visual neuron builder (drag-and-drop)
- Inline help/tooltips
- Example neurons to learn from
- Test environment before production

**UI Requirements:**
- Wizard-based neuron creation
- Interactive tutorials
- Sandbox mode (dry-run everything)
- Visual diff for YAML changes

### Persona 3: "Jordan - Platform Architect"

**Profile:**
- 12 years experience, designs infrastructure
- Manages edge IoT fleet (200+ devices)
- Needs high-level visibility + drill-down

**Needs:**
- Fleet-wide health dashboard
- Aggregated metrics from all nodes
- Anomaly detection (outlier neurons)
- Cost tracking (API usage for AI)

**UI Requirements:**
- Multi-cluster view
- Custom dashboards
- Prometheus/Grafana integration
- CSV/JSON export for analysis

---

## Architecture

### High-Level System Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         User's Browser                           │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Dashboard   │  │   Neuron     │  │    Logs      │          │
│  │    View      │  │   Builder    │  │   Viewer     │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                           │                                       │
│                           ▼                                       │
│                   ┌──────────────────┐                           │
│                   │  React Frontend  │                           │
│                   │  (PWA, Offline)  │                           │
│                   └──────────────────┘                           │
└────────────────────────────┬────────────────────────────────────┘
                             │ WebSocket + REST
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Cortex Web Server (Go)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   REST API   │  │   WebSocket  │  │     Auth     │          │
│  │  (neuron CRUD)│  │  (live logs) │  │  (JWT/OAuth) │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                           │                                       │
│                           ▼                                       │
│                   ┌──────────────────┐                           │
│                   │  Event Stream    │                           │
│                   │  (Server-Sent    │                           │
│                   │   Events)        │                           │
│                   └──────────────────┘                           │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                    Cortex Core Engine                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │  Neuron      │  │   Synapse    │  │   Execution  │          │
│  │  Manager     │  │   Executor   │  │   History    │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
                    ┌──────────────────┐
                    │   Persistent     │
                    │   Storage        │
                    │ ┌──────────────┐ │
                    │ │   SQLite     │ │  (local mode)
                    │ │   PostgreSQL │ │  (k8s mode)
                    │ │   Execution  │ │
                    │ │   Logs       │ │
                    │ └──────────────┘ │
                    └──────────────────┘
```

### Component Architecture

```
cortex/
├── cmd/
│   ├── cortex/           # Main CLI binary
│   └── cortex-web/       # Web server binary
├── pkg/
│   ├── web/
│   │   ├── server/       # HTTP/WebSocket server
│   │   ├── auth/         # Authentication/Authorization
│   │   ├── api/          # REST API handlers
│   │   └── realtime/     # WebSocket/SSE manager
│   ├── neuron/           # Core neuron engine
│   ├── synapse/          # Synapse executor
│   └── storage/          # Persistence layer
└── web/                  # Frontend (React/TypeScript)
    ├── src/
    │   ├── components/   # Reusable UI components
    │   ├── pages/        # Page layouts
    │   ├── hooks/        # React hooks
    │   ├── services/     # API clients
    │   └── stores/       # State management (Zustand)
    └── public/           # Static assets
```

---

## Feature Specifications

### Phase 1: Core Features (MVP)

#### 1. Dashboard View

**Purpose:** High-level overview of neuron execution health

**Layout:**
```
┌───────────────────────────────────────────────────────────────┐
│  Cortex Dashboard                    [Dark Mode] [Profile▼]   │
├───────────────────────────────────────────────────────────────┤
│  Time Range: [Last 24h ▼]    Refresh: [Auto ✓] [Manual ↻]   │
├────────────────────┬──────────────────┬───────────────────────┤
│  Total Executions  │  Success Rate    │  Active Synapses      │
│                    │                  │                       │
│      1,234         │     98.2%        │        5              │
│  ↑ 15% vs prev day │  ↓ 0.3% vs prev  │   ↑ 2 new today       │
└────────────────────┴──────────────────┴───────────────────────┘

┌───────────────────────────────────────────────────────────────┐
│  Execution Timeline                                            │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │         Success █████████████████████████ 95%            │ │
│  │         Warning ██ 3%                                    │ │
│  │         Failed  █ 2%                                     │ │
│  └──────────────────────────────────────────────────────────┘ │
│    6am     9am     12pm     3pm     6pm     9pm    12am       │
└───────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────┐
│  Recent Executions                       [View All→]           │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ ✓ check_disk_space      2 min ago      Exit 0  [Logs]   │ │
│  │ ✓ check_nginx_status    5 min ago      Exit 0  [Logs]   │ │
│  │ ⚠ check_ssl_cert        8 min ago      Exit 120[Logs]   │ │
│  │ ✗ restart_apache        10 min ago     Exit 1  [Logs]   │ │
│  │ ✓ backup_database       15 min ago     Exit 0  [Logs]   │ │
│  └──────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────┘

┌───────────────────────────────────────────────────────────────┐
│  Top Failing Neurons                    [Fix Suggestions→]    │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ 1. check_api_health      12 failures    Exit 120, 121   │ │
│  │    💡 AI: "API timeout - increase threshold to 10s"      │ │
│  │ 2. check_db_replication  5 failures     Exit 122         │ │
│  │    💡 AI: "Replication lag - add auto-failover neuron"   │ │
│  └──────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────┘
```

**Interactions:**
- Click metric card → Drill down to filtered view
- Hover on timeline bar → Tooltip with exact counts
- Click recent execution → Open detailed log view
- Click AI suggestion → Open neuron editor with fix applied

#### 2. Neuron Library

**Purpose:** Browse, search, and manage neurons

**Layout:**
```
┌───────────────────────────────────────────────────────────────┐
│  Neuron Library                                                │
├───────────────────────────────────────────────────────────────┤
│  🔍 Search neurons...                  [+ Create Neuron]       │
│                                                                │
│  Filters: [Type▼] [Platform▼] [Tags▼]       Sort: [Recent▼]  │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│  📁 My Neurons (23)                              [Collapse]    │
│  ├─ ✓ check_disk_space          Linux  Check   ⭐42  👤me    │
│  │  │  Last run: 2 min ago (success)     [Run] [Edit] [⋮]    │
│  │                                                             │
│  ├─ ✓ check_nginx_status        Linux  Check   ⭐18  👤me    │
│  │  │  Last run: 5 min ago (success)     [Run] [Edit] [⋮]    │
│  │                                                             │
│  └─ ⚠ check_ssl_cert            Linux  Check   ⭐8   👤me    │
│     │  Last run: 8 min ago (warning)     [Run] [Edit] [⋮]    │
│                                                                │
│  📁 Community Neurons (156)                      [Collapse]    │
│  ├─ ✓ k8s_pod_health_check      K8s    Check   ⭐234 👤john │
│  │  │  Monitor Kubernetes pod status    [Install] [Preview]  │
│  │                                                             │
│  ├─ ✓ postgres_replication_lag  DB     Check   ⭐189 👤sarah│
│  │  │  Check PostgreSQL lag < 5s        [Install] [Preview]  │
│  │                                                             │
│  └─ ⚠ aws_cost_alert            Cloud  Check   ⭐156 👤alex │
│     │  Alert on AWS spend > $X          [Install] [Preview]  │
│                                                                │
│  📁 Team Neurons (8)                             [Collapse]    │
│  └─ ✓ production_health_suite   Multi  Synapse ⭐12  👤team │
│     │  Full prod health check           [Run] [View]         │
└───────────────────────────────────────────────────────────────┘
```

**Interactions:**
- Search: Real-time fuzzy search across names, descriptions, tags
- Filter: Multi-select filters (type, platform, author, status)
- Star/favorite: Bookmark frequently used neurons
- Quick run: Execute neuron directly from list
- Drag-to-synapse: Drag neuron into synapse builder

#### 3. AI Neuron Generator (Visual)

**Purpose:** Generate neurons using natural language (visual wrapper for CLI feature)

**Layout:**
```
┌───────────────────────────────────────────────────────────────┐
│  ✨ AI Neuron Generator                          [Close ✕]    │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│  What should this neuron do?                                  │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ Check if PostgreSQL replication lag is under 5 seconds   │ │
│  └──────────────────────────────────────────────────────────┘ │
│  💡 Examples: "Monitor disk space", "Restart nginx service"   │
│                                                                │
│  Advanced Options (optional)                [Show/Hide]       │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ Type:        ● Check  ○ Mutate                           │ │
│  │ Platform:    ● Linux  ○ Windows  ○ Auto-detect           │ │
│  │ Language:    ● Bash   ○ PowerShell  ○ Python             │ │
│  │ Generate tests: ✓ Yes                                    │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│                      [Generate with AI]                        │
│                                                                │
├───────────────────────────────────────────────────────────────┤
│  🤖 Generating neuron...                         [Cancel]     │
│  ████████████████░░░░░░░░░  75%                               │
│  ✓ Analyzing description                                      │
│  ✓ Finding similar neurons (found 3)                          │
│  ⏳ Generating code with GPT-4...                             │
│  ⏸ Validating YAML...                                         │
│  ⏸ Creating tests...                                          │
└───────────────────────────────────────────────────────────────┘

[After generation completes]

┌───────────────────────────────────────────────────────────────┐
│  ✅ Neuron Generated Successfully!                             │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│  📄 check_postgres_replication_lag                            │
│                                                                │
│  ┌────────────┬─────────────────────────────────────────────┐ │
│  │ neuron.yaml│                                              │ │
│  │            │ ---                                          │ │
│  │ run.sh     │ name: check_postgres_replication_lag        │ │
│  │            │ type: check                                 │ │
│  │ README.md  │ description: "Check if PostgreSQL..."        │ │
│  │            │ exec_file: run.sh                           │ │
│  │ run_test.sh│ assertExitStatus: [0]                       │ │
│  │            │ post_exec_fail_debug:                       │ │
│  └────────────┤   120: "Replication lag 5-10s"              │ │
│                │   121: "Replication lag > 10s"              │ │
│                │                                              │ │
│                │ [Switch to run.sh] [README.md] [run_test.sh]│ │
│                └──────────────────────────────────────────────┘ │
│                                                                │
│  Exit Codes:                                                  │
│  • 0: Success - Replication lag healthy (< 5 seconds)         │
│  • 120: Warning - Replication lag 5-10 seconds                │
│  • 121: Critical - Replication lag > 10 seconds               │
│  • 122: Error - Cannot connect to PostgreSQL                  │
│                                                                │
│  [Test Neuron] [Save to Library] [Edit Manually]             │
│                                                                │
│  💡 AI Suggestion: Add this neuron to your "database_health"  │
│     synapse with auto-fix for exit code 121                   │
│                      [Add to Synapse]                          │
└───────────────────────────────────────────────────────────────┘
```

**Interactions:**
- Real-time validation: Check description as user types
- Context loading: Show similar neurons while generating
- Live preview: Switch between generated files
- Inline editing: Modify generated code before saving
- One-click test: Run neuron in sandbox environment

#### 4. Synapse Builder (Visual DAG)

**Purpose:** Create and visualize neuron orchestration workflows

**Layout:**
```
┌───────────────────────────────────────────────────────────────┐
│  Synapse Builder: production_health_check      [Save] [Run]   │
├───────────────────────────────────────────────────────────────┤
│  Toolbar: [➕ Add Neuron] [🔗 Connect] [⚙️ Settings] [💾]     │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│        ┌─────────────────────┐                                │
│        │   START             │                                │
│        └──────────┬──────────┘                                │
│                   │                                            │
│         ┌─────────┴─────────┐                                 │
│         │                   │                                  │
│    ┌────▼────┐         ┌───▼─────┐                           │
│    │ check_  │         │ check_  │                           │
│    │ disk    │         │ nginx   │                           │
│    │ (2s)    │         │ (1s)    │                           │
│    └────┬────┘         └───┬─────┘                           │
│         │                  │                                  │
│         └─────────┬────────┘                                  │
│                   │                                            │
│            ┌──────▼──────┐                                    │
│            │ check_db    │                                    │
│            │ replication │                                    │
│            │ (3s)        │                                    │
│            └──────┬──────┘                                    │
│                   │                                            │
│             [Exit Code?]                                       │
│            ┌──────┼──────┐                                    │
│            │      │      │                                     │
│      Exit 0│      │120   │121                                 │
│            │      │      │                                     │
│      ┌─────▼┐  ┌─▼────┐ ┌▼─────────┐                         │
│      │ END  │  │alert_│ │failover_ │                         │
│      │(✓)   │  │team  │ │database  │                         │
│      └──────┘  └─┬────┘ └┬─────────┘                         │
│                  │        │                                    │
│                  └────┬───┘                                    │
│                       │                                        │
│                  ┌────▼────┐                                   │
│                  │  END    │                                   │
│                  │  (⚠)    │                                   │
│                  └─────────┘                                   │
│                                                                │
│  [Neuron Library Panel]                                       │
│  Drag neurons here:                                           │
│  • check_api_health                                           │
│  • rotate_logs                                                │
│  • backup_database                                            │
└───────────────────────────────────────────────────────────────┘

[Right Panel - Selected Neuron Properties]
┌───────────────────────────────────────┐
│ Properties: check_db_replication      │
├───────────────────────────────────────┤
│ Name: check_db_replication            │
│ Type: Check                           │
│ Timeout: 30s [Edit]                   │
│ Retries: 3 [Edit]                     │
│                                       │
│ On Failure:                           │
│  Exit 120 → [alert_team      ▼]      │
│  Exit 121 → [failover_database ▼]    │
│  Exit 122 → [notify_oncall    ▼]     │
│                                       │
│ [Test This Neuron]                    │
│ [View Code]                           │
└───────────────────────────────────────┘
```

**Interactions:**
- Drag-and-drop: Add neurons from library
- Auto-layout: Suggest optimal DAG arrangement
- Connection drawing: Click and drag to connect nodes
- Parallel branches: Automatic visual spacing
- Exit code routing: Visual conditional paths
- Live validation: Detect circular dependencies
- Zoom/pan: Navigate large workflows
- Export: Generate synapse.yaml or CLI command

#### 5. Execution Logs (Real-Time)

**Purpose:** Monitor neuron execution with live streaming logs

**Layout:**
```
┌───────────────────────────────────────────────────────────────┐
│  Execution: check_disk_space #12345         [Stop] [Re-run]   │
├───────────────────────────────────────────────────────────────┤
│  Status: ⏳ Running     Duration: 00:00:03    Exit Code: -    │
│  Started: 2025-01-07 14:23:45 UTC            By: sarah@dev    │
├───────────────────────────────────────────────────────────────┤
│  Logs (Live)                       [⬇️ Download] [🔍 Search]   │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ 14:23:45 [INFO] Starting neuron check_disk_space         │ │
│  │ 14:23:45 [INFO] Checking disk usage on /dev/sda1         │ │
│  │ 14:23:46 [INFO] Current usage: 75% (150GB / 200GB)       │ │
│  │ 14:23:46 [WARN] Usage approaching threshold (80%)        │ │
│  │ 14:23:47 [INFO] Checking /dev/sda2...                    │ │
│  │ 14:23:47 [INFO] Current usage: 45% (90GB / 200GB)        │ │
│  │ 14:23:48 [INFO] ✓ All disks within threshold             │ │
│  │ 14:23:48 [INFO] Neuron completed successfully            │ │
│  │ 14:23:48 [INFO] Exit code: 0                             │ │
│  │ _  ← Live cursor                                         │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  Filters: [Level: All ▼] [Search: ________]  [Auto-scroll ✓] │
│                                                                │
│  Metadata:                                                    │
│  • Execution ID: exec-12345-abc                               │
│  • Host: edge-node-42.example.com                            │
│  • Cortex Version: v1.0.0                                     │
│  • Neuron Path: /neurons/check_disk_space                     │
│                                                                │
│  Output Files:                                                │
│  • stdout.log (2.3 KB) [Download]                            │
│  • stderr.log (0 bytes) [Download]                           │
│  • execution_metadata.json (1.1 KB) [Download]               │
└───────────────────────────────────────────────────────────────┘
```

**Interactions:**
- Live streaming: WebSocket updates as logs arrive
- Auto-scroll: Follow new logs (toggle on/off)
- Level filtering: Show only ERROR/WARN/INFO/DEBUG
- Text search: Highlight matching lines
- Line numbers: Click to copy/share specific line
- Download: Export full log file
- Share URL: Link to specific execution with auth

#### 6. Fleet View (Multi-Node)

**Purpose:** Monitor Cortex instances across edge devices/clusters

**Layout:**
```
┌───────────────────────────────────────────────────────────────┐
│  Fleet Management                           [+ Add Instance]  │
├───────────────────────────────────────────────────────────────┤
│  Filters: [Status▼] [Region▼] [Tag▼]      View: [Grid] [Map] │
├───────────────────────────────────────────────────────────────┤
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ 🟢 Online    │  │ 🟡 Warning   │  │ 🔴 Offline   │        │
│  │    42        │  │     3        │  │     1        │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                │
│  Instances:                                                   │
│  ┌──────────────────────────────────────────────────────────┐ │
│  │ 🟢 edge-node-01  │ us-west-2 │ 42 exec │ Uptime: 99.8%   │ │
│  │    Raspberry Pi 4│ 4GB RAM   │ Last: 2m ago  [Manage→]   │ │
│  ├──────────────────────────────────────────────────────────┤ │
│  │ 🟢 edge-node-02  │ us-west-2 │ 38 exec │ Uptime: 99.9%   │ │
│  │    Intel NUC     │ 8GB RAM   │ Last: 1m ago  [Manage→]   │ │
│  ├──────────────────────────────────────────────────────────┤ │
│  │ 🟡 edge-node-03  │ eu-west-1 │ 51 exec │ Uptime: 95.2%   │ │
│  │    Raspberry Pi 3│ 2GB RAM   │ Last: 15m ago [Manage→]   │ │
│  │    ⚠ High memory usage (85%)                             │ │
│  ├──────────────────────────────────────────────────────────┤ │
│  │ 🔴 edge-node-04  │ ap-south-1│ 0 exec  │ Uptime: 0%      │ │
│  │    OFFLINE       │           │ Last: 3h ago  [Manage→]   │ │
│  │    ⚠ No heartbeat received                               │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                │
│  Bulk Actions: [Select All] [Deploy Neuron] [Update Version] │
└───────────────────────────────────────────────────────────────┘
```

**Interactions:**
- Real-time status: WebSocket heartbeats from each node
- Filtering: By region, status, tag, hardware type
- Bulk operations: Deploy neurons to multiple nodes
- Drill-down: Click instance → See detailed metrics
- Map view: Geo-distributed instance visualization
- Alerts: Notification when node goes offline

---

## UI/UX Design Patterns

### Design System

**Colors (Dark Theme Primary):**
```css
--bg-primary: #1a1a1a;
--bg-secondary: #2d2d2d;
--bg-tertiary: #3a3a3a;

--text-primary: #ffffff;
--text-secondary: #b0b0b0;
--text-tertiary: #808080;

--accent-primary: #00d4ff;    /* Cortex blue */
--accent-secondary: #7c3aed;  /* Purple for AI features */

--success: #10b981;
--warning: #f59e0b;
--error: #ef4444;
--info: #3b82f6;
```

**Typography:**
```css
--font-sans: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
--font-mono: 'JetBrains Mono', 'Fira Code', monospace;

--text-xs: 0.75rem;   /* 12px */
--text-sm: 0.875rem;  /* 14px */
--text-base: 1rem;    /* 16px */
--text-lg: 1.125rem;  /* 18px */
--text-xl: 1.25rem;   /* 20px */
--text-2xl: 1.5rem;   /* 24px */
```

**Spacing (8px grid):**
```css
--spacing-1: 0.25rem;  /* 4px */
--spacing-2: 0.5rem;   /* 8px */
--spacing-3: 0.75rem;  /* 12px */
--spacing-4: 1rem;     /* 16px */
--spacing-6: 1.5rem;   /* 24px */
--spacing-8: 2rem;     /* 32px */
```

### Component Library

**Buttons:**
```jsx
// Primary action
<Button variant="primary">Generate Neuron</Button>

// Secondary action
<Button variant="secondary">Cancel</Button>

// Destructive action
<Button variant="danger">Delete</Button>

// Icon button
<Button variant="ghost" icon={<PlayIcon />}>Run</Button>
```

**Status Indicators:**
```jsx
// Success
<Badge variant="success">Running</Badge>

// Warning
<Badge variant="warning">Degraded</Badge>

// Error
<Badge variant="error">Failed</Badge>

// Info
<Badge variant="info">Scheduled</Badge>
```

**Loading States:**
```jsx
// Skeleton loader for cards
<CardSkeleton />

// Spinner for actions
<Spinner size="sm" />

// Progress bar for long operations
<ProgressBar value={75} max={100} />
```

### Accessibility Features

**Keyboard Navigation:**
- All actions accessible via keyboard
- Focus indicators clearly visible
- Keyboard shortcuts displayed in tooltips
- Tab order follows logical flow

**Screen Reader Support:**
- ARIA labels on all interactive elements
- Live regions for real-time updates
- Semantic HTML structure
- Alt text for all visualizations

**Color Contrast:**
- WCAG AAA compliance for text (7:1 ratio)
- Non-color indicators for status (icons + color)
- High contrast mode toggle

**Responsive Design:**
- Mobile-first approach
- Breakpoints: 640px (sm), 768px (md), 1024px (lg), 1280px (xl)
- Touch-friendly targets (min 44x44px)
- Collapsible sidebars on mobile

---

## Kubernetes Deployment

### Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                        │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Namespace: cortex                                    │   │
│  │                                                        │   │
│  │  ┌────────────────┐      ┌────────────────┐         │   │
│  │  │   Ingress      │──────│  LoadBalancer  │         │   │
│  │  │  (nginx/traefik)│      │  (Cloud LB)    │         │   │
│  │  └────────┬───────┘      └────────────────┘         │   │
│  │           │                                           │   │
│  │  ┌────────▼────────────────────────────────────┐    │   │
│  │  │  Service: cortex-web (ClusterIP)            │    │   │
│  │  └────────┬────────────────────────────────────┘    │   │
│  │           │                                           │   │
│  │  ┌────────▼──────────────────────────────────┐      │   │
│  │  │  Deployment: cortex-web                    │      │   │
│  │  │  ┌──────────────┐  ┌──────────────┐       │      │   │
│  │  │  │   Pod 1      │  │   Pod 2      │       │      │   │
│  │  │  │              │  │              │       │      │   │
│  │  │  │  - cortex-web│  │  - cortex-web│       │      │   │
│  │  │  │    (Go API)  │  │    (Go API)  │       │      │   │
│  │  │  │              │  │              │       │      │   │
│  │  │  │  - frontend  │  │  - frontend  │       │      │   │
│  │  │  │    (nginx)   │  │    (nginx)   │       │      │   │
│  │  │  └──────┬───────┘  └──────┬───────┘       │      │   │
│  │  └─────────┼──────────────────┼──────────────┘      │   │
│  │            │                  │                       │   │
│  │  ┌─────────▼──────────────────▼─────────────┐       │   │
│  │  │  Service: cortex-postgres (ClusterIP)    │       │   │
│  │  └────────┬─────────────────────────────────┘       │   │
│  │           │                                           │   │
│  │  ┌────────▼────────────────────────────────────┐    │   │
│  │  │  StatefulSet: cortex-postgres               │    │   │
│  │  │  ┌──────────────┐                           │    │   │
│  │  │  │ PostgreSQL   │                           │    │   │
│  │  │  │ (Primary)    │                           │    │   │
│  │  │  └──────┬───────┘                           │    │   │
│  │  │         │                                    │    │   │
│  │  │  ┌──────▼────────────┐                      │    │   │
│  │  │  │ PersistentVolume  │                      │    │   │
│  │  │  │ (execution logs,  │                      │    │   │
│  │  │  │  neuron metadata) │                      │    │   │
│  │  │  └───────────────────┘                      │    │   │
│  │  └─────────────────────────────────────────────┘    │   │
│  │                                                       │   │
│  │  ┌─────────────────────────────────────────────┐    │   │
│  │  │  ConfigMap: cortex-config                    │    │   │
│  │  │  - ai.yaml                                   │    │   │
│  │  │  - server.yaml                               │    │   │
│  │  └─────────────────────────────────────────────┘    │   │
│  │                                                       │   │
│  │  ┌─────────────────────────────────────────────┐    │   │
│  │  │  Secret: cortex-secrets                      │    │   │
│  │  │  - openai-api-key                            │    │   │
│  │  │  - postgres-password                         │    │   │
│  │  │  - jwt-secret                                │    │   │
│  │  └─────────────────────────────────────────────┘    │   │
│  └───────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

### Deployment Manifests

#### 1. Namespace & ConfigMap

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cortex
  labels:
    app: cortex
    environment: production

---
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cortex-config
  namespace: cortex
data:
  server.yaml: |
    server:
      port: 8080
      host: 0.0.0.0
      tls:
        enabled: false  # Terminated at ingress
      cors:
        enabled: true
        origins:
          - https://cortex.example.com

      auth:
        enabled: true
        provider: jwt
        session_timeout: 24h

      database:
        host: cortex-postgres
        port: 5432
        database: cortex
        pool_size: 20

      storage:
        type: postgres
        retention_days: 90

  ai.yaml: |
    ai:
      default_provider: openai
      providers:
        openai:
          model: gpt-4-turbo
          timeout: 30s
        anthropic:
          model: claude-3-5-sonnet-20250107
          timeout: 30s

      generation:
        temperature: 0.2
        max_tokens: 2048
        include_similar_neurons: true
```

#### 2. Secrets

```yaml
# k8s/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: cortex-secrets
  namespace: cortex
type: Opaque
stringData:
  openai-api-key: "sk-proj-..."  # Replace with actual key
  anthropic-api-key: "sk-ant-..."
  postgres-password: "changeme"  # Generate strong password
  jwt-secret: "generate-random-256bit-key"
  admin-password: "changeme"     # Initial admin password
```

#### 3. PostgreSQL StatefulSet

```yaml
# k8s/postgres.yaml
apiVersion: v1
kind: Service
metadata:
  name: cortex-postgres
  namespace: cortex
spec:
  ports:
    - port: 5432
      targetPort: 5432
  selector:
    app: cortex-postgres
  clusterIP: None  # Headless service

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: cortex-postgres
  namespace: cortex
spec:
  serviceName: cortex-postgres
  replicas: 1
  selector:
    matchLabels:
      app: cortex-postgres
  template:
    metadata:
      labels:
        app: cortex-postgres
    spec:
      containers:
        - name: postgres
          image: postgres:16-alpine
          ports:
            - containerPort: 5432
          env:
            - name: POSTGRES_DB
              value: cortex
            - name: POSTGRES_USER
              value: cortex
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: cortex-secrets
                  key: postgres-password
            - name: PGDATA
              value: /var/lib/postgresql/data/pgdata
          volumeMounts:
            - name: postgres-storage
              mountPath: /var/lib/postgresql/data
          resources:
            requests:
              memory: "256Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "1000m"
          livenessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - cortex
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            exec:
              command:
                - pg_isready
                - -U
                - cortex
            initialDelaySeconds: 5
            periodSeconds: 5
  volumeClaimTemplates:
    - metadata:
        name: postgres-storage
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 10Gi
```

#### 4. Cortex Web Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cortex-web
  namespace: cortex
spec:
  replicas: 2
  selector:
    matchLabels:
      app: cortex-web
  template:
    metadata:
      labels:
        app: cortex-web
    spec:
      containers:
        - name: cortex-web
          image: cortex/cortex-web:v1.0.0
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: CORTEX_CONFIG
              value: /etc/cortex/server.yaml
            - name: CORTEX_AI_CONFIG
              value: /etc/cortex/ai.yaml
            - name: OPENAI_API_KEY
              valueFrom:
                secretKeyRef:
                  name: cortex-secrets
                  key: openai-api-key
            - name: ANTHROPIC_API_KEY
              valueFrom:
                secretKeyRef:
                  name: cortex-secrets
                  key: anthropic-api-key
            - name: POSTGRES_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: cortex-secrets
                  key: postgres-password
            - name: JWT_SECRET
              valueFrom:
                secretKeyRef:
                  name: cortex-secrets
                  key: jwt-secret
          volumeMounts:
            - name: config
              mountPath: /etc/cortex
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 5
      volumes:
        - name: config
          configMap:
            name: cortex-config

---
apiVersion: v1
kind: Service
metadata:
  name: cortex-web
  namespace: cortex
spec:
  type: ClusterIP
  ports:
    - port: 80
      targetPort: 8080
      protocol: TCP
  selector:
    app: cortex-web
```

#### 5. Ingress

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cortex-ingress
  namespace: cortex
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/websocket-services: "cortex-web"
    nginx.ingress.kubernetes.io/proxy-connect-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"
    nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - cortex.example.com
      secretName: cortex-tls
  rules:
    - host: cortex.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: cortex-web
                port:
                  number: 80
```

### Helm Chart

```yaml
# helm/cortex/Chart.yaml
apiVersion: v2
name: cortex
description: AI-powered infrastructure debugging orchestrator
type: application
version: 1.0.0
appVersion: "1.0.0"

---
# helm/cortex/values.yaml
replicaCount: 2

image:
  repository: cortex/cortex-web
  tag: v1.0.0
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: cortex.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: cortex-tls
      hosts:
        - cortex.example.com

postgresql:
  enabled: true
  auth:
    username: cortex
    database: cortex
    existingSecret: cortex-secrets
    secretKeys:
      userPasswordKey: postgres-password
  primary:
    persistence:
      enabled: true
      size: 10Gi
    resources:
      requests:
        memory: 256Mi
        cpu: 250m
      limits:
        memory: 1Gi
        cpu: 1000m

config:
  server:
    port: 8080
    auth:
      enabled: true
  ai:
    defaultProvider: openai

secrets:
  openaiApiKey: ""  # Set via --set or values override
  anthropicApiKey: ""
  jwtSecret: ""
  adminPassword: ""

resources:
  requests:
    memory: 128Mi
    cpu: 100m
  limits:
    memory: 512Mi
    cpu: 500m

autoscaling:
  enabled: false
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 80
```

### Quick Deployment Guide

```bash
# 1. Install with Helm
helm repo add cortex https://charts.cortex.dev
helm repo update

# 2. Create namespace
kubectl create namespace cortex

# 3. Create secrets
kubectl create secret generic cortex-secrets \
  --from-literal=openai-api-key=$OPENAI_API_KEY \
  --from-literal=postgres-password=$(openssl rand -base64 32) \
  --from-literal=jwt-secret=$(openssl rand -base64 32) \
  --from-literal=admin-password=$(openssl rand -base64 16) \
  -n cortex

# 4. Install Cortex
helm install cortex cortex/cortex \
  --namespace cortex \
  --set ingress.hosts[0].host=cortex.example.com \
  --set ingress.tls[0].secretName=cortex-tls \
  --set ingress.tls[0].hosts[0]=cortex.example.com

# 5. Wait for deployment
kubectl wait --for=condition=ready pod \
  -l app=cortex-web \
  -n cortex \
  --timeout=300s

# 6. Get admin password
kubectl get secret cortex-secrets \
  -n cortex \
  -o jsonpath='{.data.admin-password}' | base64 -d

# 7. Access UI
echo "https://cortex.example.com"
```

---

## Technical Stack

### Frontend

**Framework:**
- **React 18** with TypeScript
- **Vite** for build tooling (fast HMR, optimized builds)
- **React Router v6** for routing

**State Management:**
- **Zustand** - Lightweight, simpler than Redux
- **TanStack Query** - Server state management, caching

**UI Components:**
- **Tailwind CSS** - Utility-first styling
- **Radix UI** - Accessible primitives (headless)
- **Lucide React** - Icon library

**Real-time:**
- **WebSocket** - Live log streaming
- **Server-Sent Events (SSE)** - Dashboard updates

**Visualization:**
- **Recharts** - Dashboard charts
- **React Flow** - DAG editor for synapse builder

**PWA:**
- **Workbox** - Service worker generation
- **Offline storage** - IndexedDB for cached data

### Backend

**Framework:**
- **Go 1.22+** with standard library
- **Chi Router** - Lightweight HTTP router
- **gorilla/websocket** - WebSocket support

**API:**
- **REST** - CRUD operations (JSON)
- **WebSocket** - Real-time updates
- **OpenAPI 3.0** - API documentation (generated)

**Database:**
- **SQLite** - Local/development mode
- **PostgreSQL 16** - Production/Kubernetes mode
- **GORM** - ORM for database operations

**Authentication:**
- **JWT** - Stateless auth tokens
- **bcrypt** - Password hashing
- **OAuth 2.0** - SSO integration (optional)

**Observability:**
- **Prometheus** - Metrics export
- **OpenTelemetry** - Distributed tracing
- **Structured logging** - JSON logs via zerolog

---

## Implementation Roadmap

### Phase 1: MVP (Weeks 1-6)

**Week 1-2: Backend Infrastructure**
- [ ] Web server with REST API
- [ ] WebSocket server for real-time
- [ ] SQLite storage layer
- [ ] Authentication (JWT)
- [ ] OpenAPI spec generation

**Week 3-4: Frontend Foundation**
- [ ] React project setup (Vite + TS)
- [ ] Design system components
- [ ] Dashboard view
- [ ] Neuron library view
- [ ] Execution logs view

**Week 5-6: Integration & Polish**
- [ ] Connect frontend to backend
- [ ] Real-time log streaming
- [ ] AI neuron generator UI
- [ ] Basic synapse builder
- [ ] Docker container builds

### Phase 2: Kubernetes & Advanced Features (Weeks 7-12)

**Week 7-8: Kubernetes Deployment**
- [ ] PostgreSQL StatefulSet
- [ ] Helm chart creation
- [ ] Ingress configuration
- [ ] RBAC implementation
- [ ] Multi-replica support

**Week 9-10: Advanced UI**
- [ ] Visual synapse builder (DAG)
- [ ] Fleet management view
- [ ] Custom dashboards
- [ ] Mobile optimizations
- [ ] PWA capabilities

**Week 11-12: Production Readiness**
- [ ] Performance optimizations
- [ ] Security hardening
- [ ] Comprehensive testing
- [ ] Documentation
- [ ] Beta release

---

## Accessibility & Performance

### Accessibility (WCAG 2.1 AAA)

**Keyboard Navigation:**
```javascript
// All interactive elements keyboard accessible
<Button onKeyDown={handleKeyDown} tabIndex={0}>
  Run Neuron
</Button>

// Keyboard shortcuts
useEffect(() => {
  const handleGlobalKeyboard = (e) => {
    if (e.metaKey && e.key === 'k') {
      openCommandPalette(); // Cmd+K for search
    }
    if (e.metaKey && e.key === 'n') {
      createNeuron(); // Cmd+N for new neuron
    }
  };
  window.addEventListener('keydown', handleGlobalKeyboard);
}, []);
```

**Screen Reader Support:**
```jsx
// ARIA labels for context
<button aria-label="Run neuron check_disk_space">
  <PlayIcon />
</button>

// Live regions for updates
<div role="log" aria-live="polite" aria-atomic="false">
  {newLogLines.map(line => <p key={line.id}>{line.text}</p>)}
</div>

// Status announcements
<div className="sr-only" role="status">
  Neuron execution completed successfully
</div>
```

### Performance Targets

**Load Time:**
- Initial load: < 2s (3G network)
- Time to Interactive (TTI): < 3s
- First Contentful Paint (FCP): < 1s

**Bundle Size:**
- JavaScript: < 300KB gzipped
- CSS: < 50KB gzipped
- Total page weight: < 1MB

**Runtime Performance:**
- 60 FPS animations
- WebSocket latency: < 100ms
- API response time: < 200ms (p95)
- Real-time log streaming: < 50ms delay

**Optimization Techniques:**
- Code splitting by route
- Lazy loading components
- Image optimization (WebP, lazy load)
- Memoization for expensive computations
- Virtual scrolling for long lists
- Service worker caching

---

## Appendix

### A. API Endpoints

```
# Neurons
GET    /api/v1/neurons              # List neurons
POST   /api/v1/neurons              # Create neuron
GET    /api/v1/neurons/:id          # Get neuron details
PUT    /api/v1/neurons/:id          # Update neuron
DELETE /api/v1/neurons/:id          # Delete neuron
POST   /api/v1/neurons/:id/execute  # Execute neuron

# AI Generation
POST   /api/v1/ai/generate          # Generate neuron with AI
POST   /api/v1/ai/suggest           # Get AI suggestions

# Synapses
GET    /api/v1/synapses             # List synapses
POST   /api/v1/synapses             # Create synapse
GET    /api/v1/synapses/:id         # Get synapse details
PUT    /api/v1/synapses/:id         # Update synapse
DELETE /api/v1/synapses/:id         # Delete synapse
POST   /api/v1/synapses/:id/execute # Execute synapse

# Executions
GET    /api/v1/executions           # List execution history
GET    /api/v1/executions/:id       # Get execution details
GET    /api/v1/executions/:id/logs  # Get execution logs
POST   /api/v1/executions/:id/rerun # Re-run execution

# Fleet
GET    /api/v1/fleet/instances      # List all Cortex instances
GET    /api/v1/fleet/instances/:id  # Get instance details
POST   /api/v1/fleet/deploy         # Deploy neuron to instances

# WebSocket
WS     /ws/logs/:execution_id       # Stream logs
WS     /ws/dashboard                # Dashboard updates
```

### B. Database Schema

```sql
-- Neurons table
CREATE TABLE neurons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(50) NOT NULL,  -- 'check' or 'mutate'
    description TEXT,
    yaml_config TEXT NOT NULL,
    execution_script TEXT NOT NULL,
    platform VARCHAR(50),       -- 'linux', 'windows', 'darwin'
    language VARCHAR(50),       -- 'bash', 'powershell', 'python'
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    starred BOOLEAN DEFAULT FALSE,
    tags TEXT[]
);

-- Synapses table
CREATE TABLE synapses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    yaml_config TEXT NOT NULL,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Executions table
CREATE TABLE executions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    neuron_id UUID REFERENCES neurons(id),
    synapse_id UUID REFERENCES synapses(id),
    status VARCHAR(50) NOT NULL,  -- 'running', 'success', 'failed'
    exit_code INTEGER,
    started_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    duration_ms INTEGER,
    executed_by UUID REFERENCES users(id),
    instance_id UUID,  -- Fleet instance that ran this
    logs TEXT
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    role VARCHAR(50) DEFAULT 'user',  -- 'admin', 'user', 'viewer'
    created_at TIMESTAMP DEFAULT NOW(),
    last_login TIMESTAMP
);

-- Fleet instances table
CREATE TABLE fleet_instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    hostname VARCHAR(255),
    ip_address INET,
    region VARCHAR(100),
    tags TEXT[],
    status VARCHAR(50) DEFAULT 'online',  -- 'online', 'offline', 'warning'
    last_heartbeat TIMESTAMP,
    cortex_version VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

**Document Version:** 1.0
**Last Updated:** 2025-01-07
**Status:** Ready for Review
**Estimated Effort:** 12 weeks (2 engineers)
**Target Release:** Q2 2025

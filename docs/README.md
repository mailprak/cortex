# Cortex Documentation

Welcome to Cortex - an AI-powered infrastructure debugging orchestrator designed to help teams organize, share, and automate debugging workflows.

## 📚 Documentation Structure

```
docs/
├── README.md                    # This file - documentation index
├── DEPENDENCIES.md              # All dependencies and licenses
├── TESTING.md                   # Testing guide (TDD with Ginkgo/Playwright)
├── guides/                      # User and contributor guides
│   ├── getting-started.md       # Quick start guide
│   └── contributing.md          # How to contribute to Cortex
├── architecture/                # System architecture and design
│   ├── overview.md              # High-level architecture
│   ├── ai-neuron-generation.md  # AI generation system design
│   └── web-ui.md                # Web UI architecture
└── specs/                       # Technical specifications (detailed design docs)
    ├── ai-neuron-generation.md  # AI feature spec
    └── web-ui-specification.md  # Web UI spec
```

## 🚀 For New Users

**Want to use Cortex?** Start here:

1. **[Getting Started Guide](guides/getting-started.md)** - Install and run your first neuron
2. **[User Guide](guides/user-guide.md)** - Complete usage documentation
3. **[Architecture Overview](architecture/overview.md)** - Understand how Cortex works

## 🛠️ For Contributors

**Want to contribute?** Start here:

1. **[Contributing Guide](guides/contributing.md)** - How to contribute code, docs, or ideas
2. **[Testing Guide](TESTING.md)** - TDD workflow with Ginkgo and Playwright
3. **[Architecture Docs](architecture/)** - Understand the codebase
4. **[Technical Specs](specs/)** - Detailed feature specifications

## 📖 Key Documents

### User Documentation

- **[Getting Started](guides/getting-started.md)** - Install, configure, and run Cortex
- **[User Guide](guides/user-guide.md)** - Complete feature documentation
  - Creating neurons (debugging scripts)
  - Building synapses (workflows)
  - AI-powered neuron generation
  - Web UI dashboard
  - Fleet management

### Architecture Documentation

- **[Architecture Overview](architecture/overview.md)** - System design and components
- **[AI Neuron Generation](architecture/ai-neuron-generation.md)** - How AI generation works
- **[Web UI Architecture](architecture/web-ui.md)** - Frontend and backend design

### Development Documentation

- **[Contributing Guide](guides/contributing.md)** - Code contributions, issues, PRs
- **[Testing Guide](TESTING.md)** - TDD workflow, running tests, coverage
- **[Dependencies](DEPENDENCIES.md)** - All open source dependencies and licenses

### Technical Specifications

- **[AI Neuron Generation Spec](specs/ai-neuron-generation.md)** - Complete AI feature design
- **[Web UI Specification](specs/web-ui-specification.md)** - Web interface design

## 🎯 What is Cortex?

Cortex is an **open source infrastructure debugging orchestrator** that helps teams:

1. **Organize debugging knowledge** - Create reusable "neurons" (debugging scripts)
2. **Build automated workflows** - Chain neurons into "synapses" (workflows)
3. **Generate AI-powered solutions** - Describe problems in natural language, get working scripts
4. **Visualize and monitor** - Modern web UI with real-time execution tracking
5. **Deploy anywhere** - Single binary OR Kubernetes, runs on edge devices

## 🌟 Key Features

- **🤖 AI-Powered**: Generate debugging scripts from natural language
- **⚡ Edge-First**: 50MB binary, runs on Raspberry Pi
- **🎨 Modern UI**: Real-time dashboard with visual workflow builder
- **🔓 100% Open Source**: Apache 2.0, all dependencies verified open source
- **🚀 Flexible Deployment**: Local binary OR Kubernetes

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Cortex System                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│  CLI Interface          Web UI (React + Vite)       │
│  (Go Binary)            (Real-time Dashboard)       │
│       │                         │                   │
│       └─────────┬───────────────┘                   │
│                 │                                    │
│         ┌───────▼────────┐                          │
│         │  Core Engine   │                          │
│         │  - Neuron Exec │                          │
│         │  - Synapse DAG │                          │
│         │  - AI Gen      │                          │
│         └───────┬────────┘                          │
│                 │                                    │
│    ┌────────────┼────────────┐                      │
│    │            │            │                      │
│ Database    Storage      AI Providers               │
│ (SQLite/    (Local/      (OpenAI/Anthropic/         │
│  Postgres)   S3)          Ollama)                   │
└─────────────────────────────────────────────────────┘
```

See [Architecture Overview](architecture/overview.md) for details.

## 📊 Project Status

- **Current Version**: 1.0 (in development)
- **Status**: Active development
- **License**: Apache 2.0
- **Language**: Go 1.25.4
- **Dependencies**: 100% open source

### Roadmap

- ✅ Core neuron/synapse execution
- ✅ CLI interface
- 🚧 AI neuron generation (OpenAI, Anthropic, Ollama)
- 🚧 Web UI dashboard
- 🚧 Visual synapse builder
- 🚧 Kubernetes deployment
- 📋 Fleet management
- 📋 Plugin marketplace

## 🤝 Contributing

We welcome contributions! Please see:

- **[Contributing Guide](guides/contributing.md)** - How to contribute
- **[Testing Guide](TESTING.md)** - How to write and run tests
- **[GitHub Issues](https://github.com/anoop2811/cortex/issues)** - Bug reports and feature requests

### Quick Start for Contributors

```bash
# 1. Clone the repository
git clone https://github.com/anoop2811/cortex.git
cd cortex

# 2. Install dependencies
make install-deps

# 3. Run tests
make test-all

# 4. Build the binary
make build

# 5. Run Cortex
./cortex --help
```

## 📜 License

Apache License 2.0 - See [LICENSE](../LICENSE) for details.

All dependencies are verified as 100% open source. See [DEPENDENCIES.md](DEPENDENCIES.md) for the complete list.

## 🔗 Links

- **GitHub**: https://github.com/anoop2811/cortex
- **Issues**: https://github.com/anoop2811/cortex/issues
- **Discussions**: https://github.com/anoop2811/cortex/discussions

## 📞 Community

- **GitHub Discussions**: Ask questions, share ideas
- **GitHub Issues**: Report bugs, request features
- **Pull Requests**: Contribute code and documentation

---

**Maintained By**: Cortex Community
**Last Updated**: November 2025

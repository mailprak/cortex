# Cortex UI Backend Server Architecture

## Overview

This document outlines the architecture for the Cortex UI backend web server, providing REST API endpoints and WebSocket support for real-time communication with the frontend.

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Cortex UI Backend                        │
├─────────────────────────────────────────────────────────────┤
│                                                               │
│  ┌──────────────┐         ┌──────────────┐                  │
│  │   CLI Layer  │────────▶│  HTTP Server │                  │
│  │  (cmd/ui.go) │         │   (net/http) │                  │
│  └──────────────┘         └──────┬───────┘                  │
│                                   │                           │
│         ┌─────────────────────────┼─────────────────┐        │
│         │                         │                 │        │
│         ▼                         ▼                 ▼        │
│  ┌─────────────┐         ┌──────────────┐   ┌─────────────┐│
│  │   REST API  │         │  WebSocket   │   │   Static    ││
│  │  Endpoints  │         │   Handler    │   │   Server    ││
│  └──────┬──────┘         └──────┬───────┘   └─────────────┘│
│         │                       │                            │
│         └───────────┬───────────┘                            │
│                     ▼                                         │
│         ┌────────────────────────┐                           │
│         │   Business Logic Layer  │                           │
│         │  - NeuronService        │                           │
│         │  - SynapseService       │                           │
│         │  - ExecutionService     │                           │
│         │  - MetricsService       │                           │
│         └────────────────────────┘                           │
│                     │                                         │
│                     ▼                                         │
│         ┌────────────────────────┐                           │
│         │   Integration Layer     │                           │
│         │  - Existing exec logic  │                           │
│         │  - Logger package       │                           │
│         │  - File system access   │                           │
│         └────────────────────────┘                           │
└─────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
/Users/agopalakrishnan/workspace/oss/cortex/
├── cmd/
│   └── ui.go                          # CLI command for UI server
├── web/
│   ├── server/
│   │   ├── server.go                  # Main HTTP server
│   │   ├── router.go                  # HTTP router setup
│   │   ├── middleware/
│   │   │   ├── cors.go                # CORS middleware
│   │   │   ├── logging.go             # Request logging
│   │   │   └── recovery.go            # Panic recovery
│   │   ├── handlers/
│   │   │   ├── neurons.go             # Neuron endpoints
│   │   │   ├── synapses.go            # Synapse endpoints
│   │   │   ├── execute.go             # Execution endpoint
│   │   │   ├── metrics.go             # Metrics endpoint
│   │   │   ├── executions.go          # Execution history
│   │   │   └── websocket.go           # WebSocket handler
│   │   ├── services/
│   │   │   ├── neuron_service.go      # Neuron business logic
│   │   │   ├── synapse_service.go     # Synapse business logic
│   │   │   ├── execution_service.go   # Execution logic
│   │   │   ├── metrics_service.go     # Metrics collection
│   │   │   └── websocket_hub.go       # WebSocket hub
│   │   └── models/
│   │       ├── neuron.go              # Neuron models
│   │       ├── synapse.go             # Synapse models
│   │       ├── execution.go           # Execution models
│   │       └── metrics.go             # Metrics models
│   └── build/                          # Frontend build files
└── docs/
    └── ui-backend-architecture.md      # This file
```

## Component Details

### 1. CLI Layer (`cmd/ui.go`)

**Purpose**: Entry point for the UI server command

**Responsibilities**:
- Parse CLI flags (port, host, debug mode)
- Initialize server configuration
- Start HTTP server
- Handle graceful shutdown

**Example Command**:
```bash
cortex ui --port 8080 --host 0.0.0.0 --debug
```

**Flags**:
- `--port`: Server port (default: 8080)
- `--host`: Server host (default: localhost)
- `--debug`: Enable debug logging

### 2. HTTP Server (`web/server/server.go`)

**Purpose**: Core HTTP server implementation

**Key Features**:
- Graceful shutdown support
- Context-based lifecycle management
- Configurable timeouts
- TLS support (optional)

**Configuration**:
```go
type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
    TLSEnabled   bool
    TLSCert      string
    TLSKey       string
}
```

### 3. Router (`web/server/router.go`)

**Purpose**: HTTP route configuration and middleware setup

**Routes**:
```
GET  /                          → Serve frontend
GET  /api/neurons               → List neurons
GET  /api/synapses              → List synapses
POST /api/execute               → Execute neuron/synapse
GET  /api/metrics               → System metrics
GET  /api/executions            → Execution history
GET  /api/executions/:id        → Specific execution
GET  /ws                        → WebSocket connection
GET  /health                    → Health check
```

### 4. Middleware

#### CORS Middleware (`middleware/cors.go`)
- Handle cross-origin requests
- Configurable allowed origins
- Support preflight requests

#### Logging Middleware (`middleware/logging.go`)
- Log all HTTP requests
- Track request duration
- Integrate with existing logger package

#### Recovery Middleware (`middleware/recovery.go`)
- Recover from panics
- Return 500 errors gracefully
- Log stack traces

### 5. API Handlers

#### Neurons Handler (`handlers/neurons.go`)

**Endpoint**: `GET /api/neurons`

**Response**:
```json
{
  "neurons": [
    {
      "name": "neuron-name",
      "path": "/path/to/neuron",
      "description": "Neuron description",
      "type": "executable|script",
      "lastModified": "2025-11-10T12:00:00Z"
    }
  ]
}
```

#### Synapses Handler (`handlers/synapses.go`)

**Endpoint**: `GET /api/synapses`

**Response**:
```json
{
  "synapses": [
    {
      "name": "synapse-name",
      "path": "/path/to/synapse",
      "neurons": ["neuron1", "neuron2"],
      "lastModified": "2025-11-10T12:00:00Z"
    }
  ]
}
```

#### Execute Handler (`handlers/execute.go`)

**Endpoint**: `POST /api/execute`

**Request**:
```json
{
  "type": "neuron|synapse",
  "name": "item-name",
  "args": ["arg1", "arg2"],
  "env": {
    "KEY": "value"
  }
}
```

**Response**:
```json
{
  "executionId": "uuid",
  "status": "running|completed|failed",
  "startTime": "2025-11-10T12:00:00Z"
}
```

#### Metrics Handler (`handlers/metrics.go`)

**Endpoint**: `GET /api/metrics`

**Response**:
```json
{
  "system": {
    "cpuUsage": 45.2,
    "memoryUsage": 1024,
    "goroutines": 23
  },
  "executions": {
    "total": 150,
    "successful": 140,
    "failed": 10,
    "running": 2
  },
  "uptime": 3600
}
```

#### Executions Handler (`handlers/executions.go`)

**Endpoint**: `GET /api/executions`

**Query Parameters**:
- `limit`: Number of results (default: 50)
- `offset`: Pagination offset
- `status`: Filter by status

**Response**:
```json
{
  "executions": [
    {
      "id": "uuid",
      "type": "neuron|synapse",
      "name": "item-name",
      "status": "completed",
      "startTime": "2025-11-10T12:00:00Z",
      "endTime": "2025-11-10T12:01:00Z",
      "duration": 60000,
      "exitCode": 0
    }
  ],
  "total": 150,
  "limit": 50,
  "offset": 0
}
```

### 6. WebSocket Handler (`handlers/websocket.go`)

**Endpoint**: `GET /ws`

**Purpose**: Real-time log streaming

**Message Types**:

**Client → Server**:
```json
{
  "type": "subscribe",
  "executionId": "uuid"
}
```

**Server → Client**:
```json
{
  "type": "log",
  "executionId": "uuid",
  "timestamp": "2025-11-10T12:00:00Z",
  "level": "info|error|debug",
  "message": "Log message"
}
```

```json
{
  "type": "status",
  "executionId": "uuid",
  "status": "running|completed|failed",
  "exitCode": 0
}
```

**Connection Management**:
- Ping/pong for connection health
- Automatic reconnection support
- Multiple subscriptions per connection

### 7. Services Layer

#### Neuron Service (`services/neuron_service.go`)

**Responsibilities**:
- Scan filesystem for neurons
- Parse neuron metadata
- Cache neuron information
- Validate neuron executability

**Methods**:
```go
type NeuronService interface {
    ListNeurons() ([]Neuron, error)
    GetNeuron(name string) (*Neuron, error)
    ValidateNeuron(name string) error
}
```

#### Synapse Service (`services/synapse_service.go`)

**Responsibilities**:
- Scan filesystem for synapses
- Parse synapse configuration
- Validate synapse dependencies

**Methods**:
```go
type SynapseService interface {
    ListSynapses() ([]Synapse, error)
    GetSynapse(name string) (*Synapse, error)
    ValidateSynapse(name string) error
}
```

#### Execution Service (`services/execution_service.go`)

**Responsibilities**:
- Execute neurons/synapses
- Track execution state
- Capture logs and output
- Manage execution history

**Methods**:
```go
type ExecutionService interface {
    Execute(req ExecuteRequest) (*Execution, error)
    GetExecution(id string) (*Execution, error)
    ListExecutions(filter ExecutionFilter) ([]Execution, error)
    StreamLogs(id string) (<-chan LogEntry, error)
}
```

#### Metrics Service (`services/metrics_service.go`)

**Responsibilities**:
- Collect system metrics
- Track execution statistics
- Calculate performance metrics

**Methods**:
```go
type MetricsService interface {
    GetSystemMetrics() (*SystemMetrics, error)
    GetExecutionMetrics() (*ExecutionMetrics, error)
    RecordExecution(execution *Execution)
}
```

#### WebSocket Hub (`services/websocket_hub.go`)

**Responsibilities**:
- Manage WebSocket connections
- Broadcast messages to subscribers
- Handle connection lifecycle

**Methods**:
```go
type WebSocketHub interface {
    Register(conn *Connection)
    Unregister(conn *Connection)
    Broadcast(executionId string, message Message)
    SubscribeToExecution(conn *Connection, executionId string)
}
```

### 8. Models

#### Neuron Model (`models/neuron.go`)
```go
type Neuron struct {
    Name         string    `json:"name"`
    Path         string    `json:"path"`
    Description  string    `json:"description"`
    Type         string    `json:"type"`
    LastModified time.Time `json:"lastModified"`
}
```

#### Synapse Model (`models/synapse.go`)
```go
type Synapse struct {
    Name         string    `json:"name"`
    Path         string    `json:"path"`
    Neurons      []string  `json:"neurons"`
    LastModified time.Time `json:"lastModified"`
}
```

#### Execution Model (`models/execution.go`)
```go
type Execution struct {
    ID        string            `json:"id"`
    Type      string            `json:"type"`
    Name      string            `json:"name"`
    Args      []string          `json:"args"`
    Env       map[string]string `json:"env"`
    Status    string            `json:"status"`
    StartTime time.Time         `json:"startTime"`
    EndTime   *time.Time        `json:"endTime,omitempty"`
    Duration  int64             `json:"duration,omitempty"`
    ExitCode  int               `json:"exitCode"`
    Logs      []LogEntry        `json:"logs,omitempty"`
}

type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    Level     string    `json:"level"`
    Message   string    `json:"message"`
}
```

#### Metrics Model (`models/metrics.go`)
```go
type SystemMetrics struct {
    CPUUsage    float64 `json:"cpuUsage"`
    MemoryUsage uint64  `json:"memoryUsage"`
    Goroutines  int     `json:"goroutines"`
}

type ExecutionMetrics struct {
    Total      int `json:"total"`
    Successful int `json:"successful"`
    Failed     int `json:"failed"`
    Running    int `json:"running"`
}

type Metrics struct {
    System     SystemMetrics     `json:"system"`
    Executions ExecutionMetrics  `json:"executions"`
    Uptime     int64            `json:"uptime"`
}
```

## Integration Points

### 1. Existing Logger Package

**Location**: `/Users/agopalakrishnan/workspace/oss/cortex/logger`

**Integration**:
- Use existing logger for all server logging
- Integrate with WebSocket for log streaming
- Maintain consistent log format

**Usage**:
```go
import "github.com/yourusername/cortex/logger"

logger.Info("Server started", "port", config.Port)
logger.Error("Execution failed", "error", err)
```

### 2. Existing Exec Logic

**Integration Strategy**:
- Import and use existing neuron/synapse execution code
- Wrap execution in service layer
- Capture output for WebSocket streaming

### 3. File System Access

**Directories to Scan**:
- Neuron directories (from config)
- Synapse directories (from config)

**Caching Strategy**:
- Cache neuron/synapse lists
- Refresh on file system events (optional)
- Manual refresh endpoint

## Security Considerations

1. **Authentication** (Future):
   - Token-based authentication
   - API key support
   - Role-based access control

2. **Input Validation**:
   - Validate all input parameters
   - Sanitize file paths
   - Prevent command injection

3. **Rate Limiting**:
   - Limit execution requests
   - Throttle WebSocket connections

4. **CORS**:
   - Configurable allowed origins
   - Secure default settings

## Error Handling

**Standard Error Response**:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable message",
    "details": {}
  }
}
```

**HTTP Status Codes**:
- 200: Success
- 400: Bad request
- 404: Not found
- 500: Internal server error
- 503: Service unavailable

## Performance Considerations

1. **Concurrency**:
   - Use goroutines for long-running executions
   - Connection pooling for WebSockets
   - Buffered channels for log streaming

2. **Caching**:
   - Cache neuron/synapse lists
   - Cache static files (if not using go:embed)

3. **Resource Limits**:
   - Maximum concurrent executions
   - Maximum WebSocket connections
   - Memory limits for log storage

## Testing Strategy

1. **Unit Tests**:
   - Test all service methods
   - Mock external dependencies
   - Test error handling

2. **Integration Tests**:
   - Test API endpoints
   - Test WebSocket communication
   - Test execution flow

3. **Load Tests**:
   - Concurrent execution stress testing
   - WebSocket connection limits
   - API endpoint performance

## Deployment

**Build Command**:
```bash
go build -o cortex-ui ./cmd/cortex
```

**Run Command**:
```bash
./cortex ui --port 8080
```

**Docker Support** (Future):
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o cortex-ui ./cmd/cortex
EXPOSE 8080
CMD ["./cortex-ui", "ui", "--port", "8080"]
```

## Future Enhancements

1. **Authentication & Authorization**
2. **Persistent Storage** (SQLite/PostgreSQL)
3. **Execution History Retention**
4. **Advanced Metrics** (Prometheus integration)
5. **Distributed Execution**
6. **Plugin System**
7. **API Versioning**
8. **GraphQL Support**

## Dependencies

**Required Packages**:
```
github.com/gorilla/websocket  - WebSocket support
github.com/gorilla/mux       - HTTP router (optional)
github.com/google/uuid       - UUID generation
```

**Standard Library**:
- net/http
- encoding/json
- context
- sync
- time

## Implementation Phases

### Phase 1: Core Infrastructure
- HTTP server setup
- Router and middleware
- Basic API endpoints (neurons, synapses)
- Static file serving

### Phase 2: Execution Support
- Execution service
- Execute endpoint
- Execution history
- Metrics collection

### Phase 3: Real-time Features
- WebSocket handler
- WebSocket hub
- Log streaming
- Status updates

### Phase 4: Polish & Testing
- Comprehensive error handling
- Unit tests
- Integration tests
- Documentation

## Conclusion

This architecture provides a solid foundation for the Cortex UI backend server with:
- Clean separation of concerns
- Scalable design
- Real-time capabilities
- Easy integration with existing codebase
- Room for future enhancements

The modular structure allows for incremental development and testing while maintaining code quality and maintainability.

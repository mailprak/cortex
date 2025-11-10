package models

import "time"

// Neuron represents a neuron definition
type Neuron struct {
	Name        string `json:"name"`
	Type        string `json:"type"`        // "check" or "mutate"
	Description string `json:"description"`
	Path        string `json:"path"`
}

// Synapse represents a synapse workflow
type Synapse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Path        string   `json:"path"`
	Neurons     []string `json:"neurons"`
}

// ExecuteRequest represents an execution request
type ExecuteRequest struct {
	Type string `json:"type"` // "neuron" or "synapse"
	Name string `json:"name"`
	Path string `json:"path"`
}

// ExecuteResponse represents an execution response
type ExecuteResponse struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"startTime"`
	Message   string    `json:"message,omitempty"`
}

// Execution represents an execution record
type Execution struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime,omitempty"`
	Duration  float64   `json:"duration,omitempty"` // seconds
	Logs      []string  `json:"logs,omitempty"`
}

// SystemMetrics represents system metrics
type SystemMetrics struct {
	CPU    float64 `json:"cpu"`    // percentage
	Memory float64 `json:"memory"` // percentage
	Disk   float64 `json:"disk"`   // percentage
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type      string      `json:"type"` // "log", "status", "metrics", "error"
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// LogMessage represents a log message
type LogMessage struct {
	ExecutionID string `json:"executionId"`
	Level       string `json:"level"` // "info", "error", "debug"
	Message     string `json:"message"`
}

// StatusMessage represents a status update
type StatusMessage struct {
	ExecutionID string `json:"executionId"`
	Status      string `json:"status"` // "running", "completed", "failed"
}

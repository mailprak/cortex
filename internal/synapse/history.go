package synapse

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ExecutionRecord represents a single execution of a synapse
type ExecutionRecord struct {
	ID            string          `json:"id"`
	SynapseName   string          `json:"synapse_name"`
	Timestamp     time.Time       `json:"timestamp"`
	Status        string          `json:"status"` // "success", "failed", "partial"
	Duration      time.Duration   `json:"duration"`
	NeuronResults []NeuronResult  `json:"neuron_results"`
	ErrorMessage  string          `json:"error_message,omitempty"`
}

// NeuronResult represents the execution result of a single neuron
type NeuronResult struct {
	Name     string        `json:"name"`
	Status   string        `json:"status"` // "success", "failed", "skipped"
	ExitCode int           `json:"exit_code"`
	Duration time.Duration `json:"duration"`
	Stdout   string        `json:"stdout"`
	Stderr   string        `json:"stderr"`
	Error    string        `json:"error,omitempty"`
}

// HistoryManager manages execution history for synapses
type HistoryManager struct {
	baseDir string
	mu      sync.RWMutex
}

// NewHistoryManager creates a new history manager with the specified base directory
func NewHistoryManager(baseDir string) *HistoryManager {
	return &HistoryManager{
		baseDir: baseDir,
	}
}

// AddExecution adds a new execution record to the history
func (hm *HistoryManager) AddExecution(synapseName string, record ExecutionRecord) error {
	if synapseName == "" {
		return errors.New("synapse name cannot be empty")
	}

	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Ensure base directory exists
	if err := os.MkdirAll(hm.baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create history directory: %w", err)
	}

	historyFile := filepath.Join(hm.baseDir, synapseName+".json")

	// Read existing history
	var history []ExecutionRecord
	if data, err := os.ReadFile(historyFile); err == nil {
		if err := json.Unmarshal(data, &history); err != nil {
			return fmt.Errorf("failed to parse existing history: %w", err)
		}
	}

	// Append new record
	history = append(history, record)

	// Write updated history
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	if err := os.WriteFile(historyFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %w", err)
	}

	return nil
}

// GetHistory retrieves all execution records for a synapse
func (hm *HistoryManager) GetHistory(synapseName string) ([]ExecutionRecord, error) {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	historyFile := filepath.Join(hm.baseDir, synapseName+".json")

	// Check if file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return []ExecutionRecord{}, nil
	}

	// Read history file
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read history file: %w", err)
	}

	var history []ExecutionRecord
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to parse history: %w", err)
	}

	return history, nil
}

// GetExecutionLogs retrieves detailed logs for a specific execution
func (hm *HistoryManager) GetExecutionLogs(synapseName, executionID string) (*ExecutionRecord, error) {
	history, err := hm.GetHistory(synapseName)
	if err != nil {
		return nil, err
	}

	// If no history exists for this synapse, return history not found
	if len(history) == 0 {
		return nil, errors.New("history not found")
	}

	// Find the specific execution
	for i := range history {
		if history[i].ID == executionID {
			return &history[i], nil
		}
	}

	// History exists but this specific execution was not found
	return nil, errors.New("execution not found")
}

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return home, nil
}

// GetDefaultHistoryDir returns the default history directory (~/.cortex/history)
func GetDefaultHistoryDir() (string, error) {
	home, err := GetHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cortex", "history"), nil
}

// NewDefaultHistoryManager creates a history manager with the default directory
func NewDefaultHistoryManager() (*HistoryManager, error) {
	historyDir, err := GetDefaultHistoryDir()
	if err != nil {
		return nil, err
	}
	return NewHistoryManager(historyDir), nil
}

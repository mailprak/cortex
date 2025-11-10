package services

import (
	"fmt"
	"os/exec"
	"sync"
	"time"

	"github.com/anoop2811/cortex/logger"
	"github.com/anoop2811/cortex/web/server/models"
	"github.com/google/uuid"
)

// ExecutionService handles execution operations
type ExecutionService struct {
	logger     *logger.StandardLogger
	executions map[string]*models.Execution
	mu         sync.RWMutex
	wsHub      *WebSocketHub
}

// NewExecutionService creates a new ExecutionService
func NewExecutionService(log *logger.StandardLogger, hub *WebSocketHub) *ExecutionService {
	return &ExecutionService{
		logger:     log,
		executions: make(map[string]*models.Execution),
		wsHub:      hub,
	}
}

// Execute executes a neuron or synapse
func (s *ExecutionService) Execute(req models.ExecuteRequest) (*models.ExecuteResponse, error) {
	executionID := uuid.New().String()
	execution := &models.Execution{
		ID:        executionID,
		Type:      req.Type,
		Name:      req.Name,
		Status:    "running",
		StartTime: time.Now(),
		Logs:      []string{},
	}

	s.mu.Lock()
	s.executions[executionID] = execution
	s.mu.Unlock()

	// Send status update via WebSocket
	s.sendWebSocketMessage("status", models.StatusMessage{
		ExecutionID: executionID,
		Status:      "running",
	})

	// Execute asynchronously
	go s.runExecution(executionID, req)

	return &models.ExecuteResponse{
		ID:        executionID,
		Status:    "running",
		StartTime: execution.StartTime,
		Message:   fmt.Sprintf("Started execution of %s: %s", req.Type, req.Name),
	}, nil
}

// runExecution runs the actual execution
func (s *ExecutionService) runExecution(executionID string, req models.ExecuteRequest) {
	s.mu.RLock()
	execution := s.executions[executionID]
	s.mu.RUnlock()

	// Send initial log
	s.sendLog(executionID, "info", fmt.Sprintf("Executing %s: %s", req.Type, req.Name))

	// Execute cortex command
	var cmd *exec.Cmd
	if req.Type == "neuron" {
		cmd = exec.Command("cortex", "exec", "-p", req.Path)
	} else {
		cmd = exec.Command("cortex", "exec", "-p", req.Path)
	}

	output, err := cmd.CombinedOutput()

	execution.EndTime = time.Now()
	execution.Duration = execution.EndTime.Sub(execution.StartTime).Seconds()

	if err != nil {
		execution.Status = "failed"
		execution.Logs = append(execution.Logs, string(output))
		s.sendLog(executionID, "error", fmt.Sprintf("Execution failed: %v", err))
		s.sendWebSocketMessage("status", models.StatusMessage{
			ExecutionID: executionID,
			Status:      "failed",
		})
		return
	}

	execution.Status = "completed"
	execution.Logs = append(execution.Logs, string(output))
	s.sendLog(executionID, "info", "Execution completed successfully")
	s.sendWebSocketMessage("status", models.StatusMessage{
		ExecutionID: executionID,
		Status:      "completed",
	})
}

// ListExecutions returns all executions
func (s *ExecutionService) ListExecutions() []models.Execution {
	s.mu.RLock()
	defer s.mu.RUnlock()

	executions := make([]models.Execution, 0, len(s.executions))
	for _, exec := range s.executions {
		executions = append(executions, *exec)
	}

	return executions
}

// sendLog sends a log message via WebSocket
func (s *ExecutionService) sendLog(executionID, level, message string) {
	s.sendWebSocketMessage("log", models.LogMessage{
		ExecutionID: executionID,
		Level:       level,
		Message:     message,
	})
}

// sendWebSocketMessage sends a message to all WebSocket clients
func (s *ExecutionService) sendWebSocketMessage(msgType string, data interface{}) {
	msg := models.WebSocketMessage{
		Type:      msgType,
		Timestamp: time.Now(),
		Data:      data,
	}

	s.wsHub.Broadcast(msg)
}

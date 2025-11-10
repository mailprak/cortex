package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"syscall"

	"github.com/anoop2811/cortex/logger"
	"github.com/anoop2811/cortex/web/server/models"
	"github.com/anoop2811/cortex/web/server/services"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Handlers holds all HTTP handlers
type Handlers struct {
	logger           *logger.StandardLogger
	neuronService    *services.NeuronService
	executionService *services.ExecutionService
	wsHub            *services.WebSocketHub
}

// NewHandlers creates a new Handlers instance
func NewHandlers(log *logger.StandardLogger) *Handlers {
	hub := services.NewWebSocketHub()
	go hub.Run()

	return &Handlers{
		logger:           log,
		neuronService:    services.NewNeuronService(log),
		executionService: services.NewExecutionService(log, hub),
		wsHub:            hub,
	}
}

// ListNeurons handles GET /api/neurons
func (h *Handlers) ListNeurons(w http.ResponseWriter, r *http.Request) {
	neurons, err := h.neuronService.ListNeurons()
	if err != nil {
		h.logger.Error(err, "Failed to list neurons")
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, neurons)
}

// ListSynapses handles GET /api/synapses
func (h *Handlers) ListSynapses(w http.ResponseWriter, r *http.Request) {
	synapses, err := h.neuronService.ListSynapses()
	if err != nil {
		h.logger.Error(err, "Failed to list synapses")
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, synapses)
}

// Execute handles POST /api/execute
func (h *Handlers) Execute(w http.ResponseWriter, r *http.Request) {
	var req models.ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	execution, err := h.executionService.Execute(req)
	if err != nil {
		h.logger.Error(err, "Execution failed")
		respondJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	respondJSON(w, http.StatusOK, execution)
}

// GetMetrics handles GET /api/metrics
func (h *Handlers) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := getSystemMetrics()
	respondJSON(w, http.StatusOK, metrics)
}

// ListExecutions handles GET /api/executions
func (h *Handlers) ListExecutions(w http.ResponseWriter, r *http.Request) {
	executions := h.executionService.ListExecutions()
	respondJSON(w, http.StatusOK, executions)
}

// WebSocketHandler handles WebSocket connections
func (h *Handlers) WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err, "WebSocket upgrade failed")
		return
	}

	client := &services.WebSocketClient{
		ID:   uuid.New().String(),
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.wsHub.RegisterClient(client)

	// Start goroutines for reading and writing
	go client.WritePump()
	go client.ReadPump(h.wsHub)
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// getSystemMetrics returns current system metrics
func getSystemMetrics() models.SystemMetrics {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	// Calculate memory usage percentage
	memPercent := float64(m.Alloc) / float64(m.Sys) * 100

	// Get disk usage (simplified)
	diskPercent := getDiskUsage()

	return models.SystemMetrics{
		CPU:    getCPUUsage(),
		Memory: memPercent,
		Disk:   diskPercent,
	}
}

// getCPUUsage returns CPU usage (simplified implementation)
func getCPUUsage() float64 {
	// This is a simplified implementation
	// In production, use a proper CPU monitoring library
	return float64(runtime.NumGoroutine()) / float64(runtime.NumCPU()) * 10
}

// getDiskUsage returns disk usage percentage
func getDiskUsage() float64 {
	// Simplified implementation using syscall.Statfs
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	if err != nil {
		return 0
	}

	total := stat.Blocks * uint64(stat.Bsize)
	free := stat.Bfree * uint64(stat.Bsize)
	used := total - free

	return float64(used) / float64(total) * 100
}

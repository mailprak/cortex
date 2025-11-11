package synapse

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/anoop2811/cortex/internal/neuron"
	log "github.com/anoop2811/cortex/logger"
	"github.com/google/uuid"
)

// Executor handles synapse execution
type Executor struct {
	logger         *log.StandardLogger
	historyManager *HistoryManager
	neuronCache    map[string]*neuron.Neuron
	environment    map[string]string
	out            io.Writer
	mu             sync.Mutex
}

// NewExecutor creates a new synapse executor
func NewExecutor(logger *log.StandardLogger, historyManager *HistoryManager, out io.Writer) *Executor {
	if out == nil {
		out = os.Stdout
	}
	return &Executor{
		logger:         logger,
		historyManager: historyManager,
		neuronCache:    make(map[string]*neuron.Neuron),
		environment:    make(map[string]string),
		out:            out,
	}
}

// SetEnvironment sets environment variables for conditional evaluation
func (e *Executor) SetEnvironment(env map[string]string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.environment = env
}

// Execute executes a synapse workflow
func (e *Executor) Execute(ctx context.Context, synapse *Synapse, synapseDir string) error {
	executionID := uuid.New().String()
	startTime := time.Now()

	e.logger.Infof("Starting synapse execution: %s (ID: %s)", synapse.Name, executionID)

	// Apply timeout if specified
	if synapse.Timeout != "" {
		timeout, err := synapse.GetTimeoutDuration()
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}
	}

	// Create execution record
	record := ExecutionRecord{
		ID:            executionID,
		SynapseName:   synapse.Name,
		Timestamp:     startTime,
		Status:        "running",
		NeuronResults: []NeuronResult{},
	}

	var executionErr error

	// Execute based on mode
	if synapse.Execution == ExecutionParallel {
		executionErr = e.executeParallel(ctx, synapse, synapseDir, &record)
	} else {
		executionErr = e.executeSequential(ctx, synapse, synapseDir, &record)
	}

	// Finalize execution record
	record.Duration = time.Since(startTime)
	if executionErr != nil {
		record.Status = "failed"
		record.ErrorMessage = executionErr.Error()
	} else {
		// Check if any neurons failed
		allSuccess := true
		for _, nr := range record.NeuronResults {
			if nr.Status == "failed" {
				allSuccess = false
				break
			}
		}
		if allSuccess {
			record.Status = "success"
		} else {
			record.Status = "partial"
		}
	}

	// Save execution history
	if e.historyManager != nil {
		if err := e.historyManager.AddExecution(synapse.Name, record); err != nil {
			e.logger.Errorf(err, "Failed to save execution history")
		}
	}

	return executionErr
}

// executeSequential executes neurons sequentially
func (e *Executor) executeSequential(ctx context.Context, synapse *Synapse, synapseDir string, record *ExecutionRecord) error {
	for _, neuronRef := range synapse.Neurons {
		select {
		case <-ctx.Done():
			return fmt.Errorf("execution timeout exceeded")
		default:
		}

		// Check condition
		if !e.evaluateCondition(neuronRef.Condition) {
			fmt.Fprintf(e.out, "Skipping: %s (condition not met)\n", neuronRef.Name)
			record.NeuronResults = append(record.NeuronResults, NeuronResult{
				Name:   neuronRef.Name,
				Status: "skipped",
			})
			continue
		}

		// Execute neuron with retry
		result := e.executeNeuronWithRetry(ctx, neuronRef, synapseDir)
		record.NeuronResults = append(record.NeuronResults, result)

		// Handle failure
		if result.Status == "failed" {
			// Execute rollback neurons if specified
			if len(neuronRef.OnFailure) > 0 {
				fmt.Fprintf(e.out, "Executing rollback for %s\n", neuronRef.Name)
				for _, rollbackNeuron := range neuronRef.OnFailure {
					e.executeNeuronByName(ctx, rollbackNeuron, synapseDir)
				}
			}

			// Stop on error if configured
			if synapse.StopOnError {
				fmt.Fprintf(e.out, "Stopping execution due to error in %s\n", neuronRef.Name)
				return fmt.Errorf("neuron %s failed: %s", neuronRef.Name, result.Error)
			}
		}
	}

	return nil
}

// executeParallel executes neurons in parallel
func (e *Executor) executeParallel(ctx context.Context, synapse *Synapse, synapseDir string, record *ExecutionRecord) error {
	fmt.Fprintf(e.out, "Executing in parallel (max concurrency: %d)\n", synapse.MaxConcurrency)

	// Build dependency graph for topological sort
	readyQueue := []NeuronRef{}
	waiting := make(map[string]NeuronRef)
	dependencies := make(map[string]map[string]bool)
	completed := make(map[string]bool)

	// Initialize
	for _, neuronRef := range synapse.Neurons {
		if len(neuronRef.DependsOn) == 0 {
			readyQueue = append(readyQueue, neuronRef)
		} else {
			waiting[neuronRef.Name] = neuronRef
			dependencies[neuronRef.Name] = make(map[string]bool)
			for _, dep := range neuronRef.DependsOn {
				dependencies[neuronRef.Name][dep] = true
			}
		}
	}

	// Concurrency control
	maxConcurrency := synapse.MaxConcurrency
	if maxConcurrency <= 0 {
		maxConcurrency = 5 // default
	}
	semaphore := make(chan struct{}, maxConcurrency)

	var wg sync.WaitGroup
	var resultsMu sync.Mutex
	results := []NeuronResult{}

	// Process neurons
	for len(readyQueue) > 0 || len(waiting) > 0 {
		// Execute ready neurons
		for len(readyQueue) > 0 {
			neuronRef := readyQueue[0]
			readyQueue = readyQueue[1:]

			wg.Add(1)
			go func(nr NeuronRef) {
				defer wg.Done()

				// Acquire semaphore
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// Check condition
				if !e.evaluateCondition(nr.Condition) {
					fmt.Fprintf(e.out, "Skipping: %s (condition not met)\n", nr.Name)
					resultsMu.Lock()
					results = append(results, NeuronResult{
						Name:   nr.Name,
						Status: "skipped",
					})
					completed[nr.Name] = true
					resultsMu.Unlock()
					return
				}

				// Execute neuron
				result := e.executeNeuronWithRetry(ctx, nr, synapseDir)

				resultsMu.Lock()
				results = append(results, result)
				completed[nr.Name] = true
				resultsMu.Unlock()

				// Handle failure
				if result.Status == "failed" && len(nr.OnFailure) > 0 {
					fmt.Fprintf(e.out, "Executing rollback for %s\n", nr.Name)
					for _, rollbackNeuron := range nr.OnFailure {
						e.executeNeuronByName(ctx, rollbackNeuron, synapseDir)
					}
				}
			}(neuronRef)
		}

		// Wait for some neurons to complete
		wg.Wait()

		// Check which waiting neurons are now ready
		resultsMu.Lock()
		for name, neuronRef := range waiting {
			allDepsComplete := true
			for dep := range dependencies[name] {
				if !completed[dep] {
					allDepsComplete = false
					break
				}
			}
			if allDepsComplete {
				readyQueue = append(readyQueue, neuronRef)
				delete(waiting, name)
			}
		}
		resultsMu.Unlock()

		// Break if no progress
		if len(readyQueue) == 0 && len(waiting) > 0 {
			return fmt.Errorf("deadlock detected: some neurons cannot execute due to unmet dependencies")
		}
	}

	// Add results to record
	record.NeuronResults = results

	return nil
}

// executeNeuronWithRetry executes a neuron with retry policy
func (e *Executor) executeNeuronWithRetry(ctx context.Context, neuronRef NeuronRef, synapseDir string) NeuronResult {
	result := NeuronResult{
		Name:   neuronRef.Name,
		Status: "success",
	}

	maxAttempts := 1
	backoff := BackoffLinear
	initialDelay := time.Second

	if neuronRef.Retry != nil {
		maxAttempts = neuronRef.Retry.MaxAttempts
		if maxAttempts < 1 {
			maxAttempts = 1
		}
		backoff = neuronRef.Retry.Backoff
		if neuronRef.Retry.InitialDelay != "" {
			if d, err := time.ParseDuration(neuronRef.Retry.InitialDelay); err == nil {
				initialDelay = d
			}
		}
	}

	var lastErr error
	startTime := time.Now()

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			result.Status = "failed"
			result.Error = "execution timeout exceeded"
			result.Duration = time.Since(startTime)
			return result
		default:
		}

		if attempt > 1 {
			delay := e.calculateBackoff(initialDelay, attempt, backoff)
			fmt.Fprintf(e.out, "Retry attempt %d/%d for %s (waiting %v)\n", attempt, maxAttempts, neuronRef.Name, delay)
			time.Sleep(delay)
		}

		fmt.Fprintf(e.out, "Executing: %s\n", neuronRef.Name)

		exitCode, stdout, stderr, err := e.executeNeuron(neuronRef.Name, synapseDir)

		result.ExitCode = exitCode
		result.Stdout = stdout
		result.Stderr = stderr

		if err == nil && exitCode == 0 {
			result.Status = "success"
			result.Duration = time.Since(startTime)
			return result
		}

		lastErr = err
	}

	// All attempts failed
	result.Status = "failed"
	result.Duration = time.Since(startTime)
	if lastErr != nil {
		result.Error = lastErr.Error()
	}

	return result
}

// executeNeuronByName executes a neuron by name (for rollback)
func (e *Executor) executeNeuronByName(ctx context.Context, name string, synapseDir string) {
	fmt.Fprintf(e.out, "Executing: %s\n", name)
	e.executeNeuron(name, synapseDir)
}

// executeNeuron executes a single neuron
func (e *Executor) executeNeuron(name string, synapseDir string) (int, string, string, error) {
	// Look for neuron in synapse directory
	neuronPath := filepath.Join(synapseDir, "neurons", name+".yml")

	// Check if neuron file exists
	if _, err := os.Stat(neuronPath); os.IsNotExist(err) {
		// Try without .yml extension
		neuronPath = filepath.Join(synapseDir, "neurons", name)
		if _, err := os.Stat(neuronPath); os.IsNotExist(err) {
			return -1, "", "", fmt.Errorf("neuron not found: %s", name)
		}
	}

	// Load neuron
	n, err := neuron.NewNeuron(e.logger, neuronPath)
	if err != nil {
		return -1, "", "", fmt.Errorf("failed to load neuron: %w", err)
	}

	// Execute neuron
	exitCode, err := n.Excite(false, e.out)

	// For now, return empty stdout/stderr as neuron.Excite doesn't capture them
	return exitCode, "", "", err
}

// evaluateCondition evaluates a conditional expression
func (e *Executor) evaluateCondition(condition string) bool {
	if condition == "" {
		return true // No condition means always execute
	}

	// Simple condition evaluation: "key == 'value'"
	// This is a placeholder - a real implementation would use a proper expression parser
	e.mu.Lock()
	defer e.mu.Unlock()

	// For now, just check if the condition string exists in environment
	// A proper implementation would parse expressions like "environment == 'staging'"
	for key, value := range e.environment {
		expected := fmt.Sprintf("%s == '%s'", key, value)
		if condition == expected {
			return true
		}
	}

	return false
}

// calculateBackoff calculates delay for retry attempts
func (e *Executor) calculateBackoff(initialDelay time.Duration, attempt int, strategy BackoffStrategy) time.Duration {
	if strategy == BackoffExponential {
		// Exponential: delay * 2^(attempt-1)
		multiplier := 1 << uint(attempt-2) // 2^(attempt-2)
		return initialDelay * time.Duration(multiplier)
	}

	// Linear: delay * attempt
	return initialDelay * time.Duration(attempt)
}

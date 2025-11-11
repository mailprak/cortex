package synapse

import (
	"fmt"
	"time"

	"github.com/anoop2811/cortex/internal/neuron"
)

// ExecutionMode defines how neurons should be executed
type ExecutionMode string

const (
	ExecutionSequential ExecutionMode = "sequential"
	ExecutionParallel   ExecutionMode = "parallel"
)

// BackoffStrategy defines retry backoff behavior
type BackoffStrategy string

const (
	BackoffExponential BackoffStrategy = "exponential"
	BackoffLinear      BackoffStrategy = "linear"
)

// Synapse represents a workflow configuration
type Synapse struct {
	Name           string              `yaml:"name"`
	Definition     []neuron.Definition `yaml:"definition"`
	Neurons        []NeuronRef         `yaml:"neurons"`
	Execution      ExecutionMode       `yaml:"execution"`
	StopOnError    bool                `yaml:"stopOnError"`
	MaxConcurrency int                 `yaml:"maxConcurrency"`
	Resources      *ResourceLimits     `yaml:"resources,omitempty"`
	Timeout        string              `yaml:"timeout,omitempty"`
}

// NeuronRef references a neuron with execution metadata
type NeuronRef struct {
	Name      string        `yaml:"name"`
	Condition string        `yaml:"condition,omitempty"`
	Retry     *RetryPolicy  `yaml:"retry,omitempty"`
	OnFailure []string      `yaml:"onFailure,omitempty"`
	DependsOn []string      `yaml:"dependsOn,omitempty"`
}

// UnmarshalYAML implements custom unmarshaling to support both string and object formats
func (nr *NeuronRef) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to unmarshal as a string first
	var name string
	if err := unmarshal(&name); err == nil {
		nr.Name = name
		return nil
	}

	// If that fails, try to unmarshal as a full object
	type neuronRefAlias NeuronRef
	var ref neuronRefAlias
	if err := unmarshal(&ref); err != nil {
		return err
	}

	*nr = NeuronRef(ref)
	return nil
}

// RetryPolicy defines retry behavior for a neuron
type RetryPolicy struct {
	MaxAttempts  int             `yaml:"maxAttempts"`
	Backoff      BackoffStrategy `yaml:"backoff"`
	InitialDelay string          `yaml:"initialDelay"`
}

// ResourceLimits defines resource constraints
type ResourceLimits struct {
	Memory string `yaml:"memory"`
}

type Plan struct {
	Config Config `yaml:"config"`
}

type Steps struct {
	Serial   []string `yaml:"serial"`
	Parallel []string `yaml:"parallel"`
}

type Config struct {
	ExitOnFirstError bool `yaml:"exit_on_first_error"`
}

// GetTimeoutDuration parses timeout string to duration
func (s *Synapse) GetTimeoutDuration() (time.Duration, error) {
	if s.Timeout == "" {
		return 0, nil
	}
	return time.ParseDuration(s.Timeout)
}

// Validate checks if the synapse configuration is valid
func (s *Synapse) Validate() error {
	if s.Name == "" {
		return fmt.Errorf("synapse name cannot be empty")
	}
	if len(s.Neurons) == 0 {
		return fmt.Errorf("synapse must have at least one neuron")
	}

	// Check for duplicate neuron names
	seen := make(map[string]bool)
	for _, neuron := range s.Neurons {
		if seen[neuron.Name] {
			return fmt.Errorf("duplicate neuron name: %s", neuron.Name)
		}
		seen[neuron.Name] = true
	}

	// Validate dependencies exist
	for _, neuron := range s.Neurons {
		for _, dep := range neuron.DependsOn {
			if !seen[dep] {
				return fmt.Errorf("neuron %s depends on non-existent neuron: %s", neuron.Name, dep)
			}
		}
	}

	// Check for circular dependencies
	if err := s.detectCircularDependencies(); err != nil {
		return err
	}

	return nil
}

// detectCircularDependencies uses DFS to detect cycles in dependency graph
func (s *Synapse) detectCircularDependencies() error {
	// Build adjacency list
	graph := make(map[string][]string)
	for _, neuron := range s.Neurons {
		graph[neuron.Name] = neuron.DependsOn
	}

	// Track visit states: 0=unvisited, 1=visiting, 2=visited
	state := make(map[string]int)

	var dfs func(string, []string) error
	dfs = func(node string, path []string) error {
		if state[node] == 1 {
			// Found a cycle
			cyclePath := append(path, node)
			return fmt.Errorf("Circular dependency detected: %v", cyclePath)
		}
		if state[node] == 2 {
			// Already visited
			return nil
		}

		state[node] = 1 // Mark as visiting
		path = append(path, node)

		for _, dep := range graph[node] {
			if err := dfs(dep, path); err != nil {
				return err
			}
		}

		state[node] = 2 // Mark as visited
		return nil
	}

	// Check each neuron
	for _, neuron := range s.Neurons {
		if state[neuron.Name] == 0 {
			if err := dfs(neuron.Name, []string{}); err != nil {
				return err
			}
		}
	}

	return nil
}

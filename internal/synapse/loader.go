package synapse

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// LoadFromDirectory loads a synapse configuration from a directory
func LoadFromDirectory(dir string) (*Synapse, error) {
	configPath := filepath.Join(dir, "config.yml")

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("synapse config not found: %s", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read synapse config: %w", err)
	}

	// Parse YAML
	var synapse Synapse
	if err := yaml.Unmarshal(data, &synapse); err != nil {
		return nil, fmt.Errorf("failed to parse synapse config: %w", err)
	}

	// Validate synapse configuration
	if err := synapse.Validate(); err != nil {
		return nil, fmt.Errorf("invalid synapse configuration: %w", err)
	}

	return &synapse, nil
}

// LoadFromFile loads a synapse configuration from a specific file
func LoadFromFile(path string) (*Synapse, error) {
	// Read config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read synapse config: %w", err)
	}

	// Parse YAML
	var synapse Synapse
	if err := yaml.Unmarshal(data, &synapse); err != nil {
		return nil, fmt.Errorf("failed to parse synapse config: %w", err)
	}

	// Validate synapse configuration
	if err := synapse.Validate(); err != nil {
		return nil, fmt.Errorf("invalid synapse configuration: %w", err)
	}

	return &synapse, nil
}

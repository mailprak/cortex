/*
Copyright © 2020 The Cortex Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/anoop2811/cortex/logger"
	"github.com/anoop2811/cortex/internal/neuron"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type SynapseConfig struct {
	Name       string             `yaml:"name"`
	Definition []NeuronDefinition `yaml:"definition"`
	Plan       PlanConfig         `yaml:"plan"`
}

type NeuronDefinition struct {
	Neuron string       `yaml:"neuron"`
	Config NeuronConfig `yaml:"config"`
}

type NeuronConfig struct {
	Path string         `yaml:"path"`
	Fix  map[int]string `yaml:"fix"`
}

type PlanConfig struct {
	Config PlanSettings `yaml:"config"`
	Steps  PlanSteps    `yaml:"steps"`
}

type PlanSettings struct {
	ExitOnFirstError bool `yaml:"exit_on_first_error"`
}

type PlanSteps struct {
	Serial   []string `yaml:"serial"`
	Parallel []string `yaml:"parallel"`
}

var synapsePath string

// fireSynapseCmd represents the fireSynapse command
var fireSynapseCmd = &cobra.Command{
	Use:     "fire-synapse",
	Short:   "Execute a synapse",
	Long:    `Execute a synapse by running all neurons in the defined plan`,
	Aliases: []string{"fs", "fire"},
	Run: func(cmd *cobra.Command, args []string) {
		fireSynapse(synapsePath)
	},
}

func init() {
	rootCmd.AddCommand(fireSynapseCmd)
	fireSynapseCmd.Flags().StringVarP(&synapsePath, "path", "p", ".", "Path to synapse directory")
}

func fireSynapse(path string) {
	logger := log.NewLogger(verbose)

	// Read synapse.yaml
	synapseFile := filepath.Join(path, "synapse.yaml")
	data, err := ioutil.ReadFile(synapseFile)
	if err != nil {
		fmt.Printf("Error reading synapse.yaml: %v\n", err)
		os.Exit(1)
	}

	var config SynapseConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Error parsing synapse.yaml: %v\n", err)
		os.Exit(1)
	}

	color.New(color.FgCyan, color.Bold).Printf("\n🧠 Firing synapse: %s\n\n", config.Name)

	// Build neuron map for easy lookup
	neuronMap := make(map[string]NeuronDefinition)
	for _, def := range config.Definition {
		neuronMap[def.Neuron] = def
	}

	exitOnFirstError := config.Plan.Config.ExitOnFirstError
	hasErrors := false

	// Execute serial neurons
	if len(config.Plan.Steps.Serial) > 0 {
		color.New(color.FgYellow).Println("▶ Executing serial neurons...")
		for _, neuronName := range config.Plan.Steps.Serial {
			if err := executeNeuron(logger, neuronName, neuronMap); err != nil {
				hasErrors = true
				if exitOnFirstError {
					color.New(color.FgRed).Printf("✗ Exiting due to error in neuron '%s'\n", neuronName)
					os.Exit(1)
				}
			}
		}
	}

	// Execute parallel neurons
	if len(config.Plan.Steps.Parallel) > 0 {
		color.New(color.FgYellow).Println("\n▶ Executing parallel neurons...")
		for _, neuronName := range config.Plan.Steps.Parallel {
			if err := executeNeuron(logger, neuronName, neuronMap); err != nil {
				hasErrors = true
				if exitOnFirstError {
					color.New(color.FgRed).Printf("✗ Exiting due to error in neuron '%s'\n", neuronName)
					os.Exit(1)
				}
			}
		}
	}

	if hasErrors {
		color.New(color.FgRed, color.Bold).Println("\n⚠ Synapse completed with errors")
		os.Exit(1)
	} else {
		color.New(color.FgGreen, color.Bold).Println("\n✓ Synapse completed successfully")
	}
}

func executeNeuron(logger *log.StandardLogger, neuronName string, neuronMap map[string]NeuronDefinition) error {
	def, exists := neuronMap[neuronName]
	if !exists {
		color.New(color.FgRed).Printf("✗ Neuron '%s' not found in definition\n", neuronName)
		return fmt.Errorf("neuron not found: %s", neuronName)
	}

	// Check if neuron path exists
	neuronConfigPath := filepath.Join(def.Config.Path, "neuron.yaml")
	if _, err := os.Stat(neuronConfigPath); os.IsNotExist(err) {
		color.New(color.FgRed).Printf("✗ Neuron config not found: %s\n", neuronConfigPath)
		return err
	}

	// Load and execute neuron
	n, err := neuron.NewNeuron(logger, neuronConfigPath)
	if err != nil {
		color.New(color.FgRed).Printf("✗ Failed to load neuron '%s': %v\n", neuronName, err)
		return err
	}

	color.New(color.FgCyan).Printf("  • %s: ", neuronName)
	exitCode, err := n.Excite(false, os.Stdout)

	// Check if exit code is in the assert_exit_status list
	isExpectedExitCode := false
	for _, expectedCode := range n.AssertExitStatus {
		if expectedCode == fmt.Sprintf("%d", exitCode) {
			isExpectedExitCode = true
			break
		}
	}

	if err != nil || (!isExpectedExitCode && exitCode != 0) {
		color.New(color.FgRed).Printf(" ✗ (exit code: %d)\n", exitCode)

		// Check if there's a fix defined for this exit code
		if fixNeuron, hasFix := def.Config.Fix[exitCode]; hasFix {
			color.New(color.FgYellow).Printf("    ↳ Attempting fix with neuron: %s\n", fixNeuron)
			if err := executeNeuron(logger, fixNeuron, neuronMap); err != nil {
				return err
			}
		}
		return fmt.Errorf("neuron failed with exit code %d", exitCode)
	}

	if exitCode != 0 {
		color.New(color.FgYellow).Printf(" ⚠ (exit code: %d)\n", exitCode)
	} else {
		color.New(color.FgGreen).Println(" ✓")
	}
	return nil
}

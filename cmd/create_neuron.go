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
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var neuronType string

// createNeuronCmd represents the createNeuron command
var createNeuronCmd = &cobra.Command{
	Use:     "create-neuron <name>",
	Short:   "Bootstrap a new neuron folder",
	Long:    `Bootstrap a new neuron folder with config and run scripts`,
	Aliases: []string{"cn"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		neuronName := args[0]
		createNeuron(neuronName)
	},
}

func init() {
	rootCmd.AddCommand(createNeuronCmd)
	createNeuronCmd.Flags().StringVarP(&neuronType, "type", "t", "check", "Type of neuron (check or mutate)")
}

func createNeuron(name string) {
	// Create neuron directory
	if err := os.MkdirAll(name, 0755); err != nil {
		fmt.Printf("Error creating neuron directory: %v\n", err)
		os.Exit(1)
	}

	// Create neuron.yaml
	neuronConfig := map[string]interface{}{
		"name":                     name,
		"type":                     neuronType,
		"description":              fmt.Sprintf("Description for %s", name),
		"exec_file":                "./run.sh",
		"pre_exec_debug":           fmt.Sprintf("Executing %s", name),
		"assert_exit_status":       []int{0},
		"post_exec_success_debug":  fmt.Sprintf("%s completed successfully", name),
		"post_exec_fail_debug": map[int]string{
			1: "Execution failed",
		},
	}

	yamlData, err := yaml.Marshal(neuronConfig)
	if err != nil {
		fmt.Printf("Error marshaling neuron config: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(name, "neuron.yaml")
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		fmt.Printf("Error writing neuron.yaml: %v\n", err)
		os.Exit(1)
	}

	// Create run.sh
	runShScript := `#!/bin/bash
# Add your debug/check logic here
echo "Running ` + name + `"
exit 0
`
	runShPath := filepath.Join(name, "run.sh")
	if err := os.WriteFile(runShPath, []byte(runShScript), 0755); err != nil {
		fmt.Printf("Error writing run.sh: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Created neuron '%s' with type '%s'\n", name, neuronType)
	fmt.Printf("  - %s/neuron.yaml\n", name)
	fmt.Printf("  - %s/run.sh\n", name)
}

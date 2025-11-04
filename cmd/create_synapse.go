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

// createSynapseCmd represents the createSynapse command
var createSynapseCmd = &cobra.Command{
	Use:     "create-synapse <name>",
	Short:   "Bootstrap a new synapse folder",
	Long:    `Bootstrap a new synapse folder with config and file structure`,
	Aliases: []string{"cs"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		synapseName := args[0]
		createSynapse(synapseName)
	},
}

func init() {
	rootCmd.AddCommand(createSynapseCmd)
}

func createSynapse(name string) {
	// Create synapse directory
	if err := os.MkdirAll(name, 0755); err != nil {
		fmt.Printf("Error creating synapse directory: %v\n", err)
		os.Exit(1)
	}

	// Create synapse.yaml
	synapseConfig := map[string]interface{}{
		"name":       name,
		"definition": []map[string]interface{}{},
		"plan": map[string]interface{}{
			"config": map[string]interface{}{
				"exit_on_first_error": false,
			},
			"steps": map[string]interface{}{
				"serial":   []string{},
				"parallel": []string{},
			},
		},
	}

	yamlData, err := yaml.Marshal(synapseConfig)
	if err != nil {
		fmt.Printf("Error marshaling synapse config: %v\n", err)
		os.Exit(1)
	}

	configPath := filepath.Join(name, "synapse.yaml")
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		fmt.Printf("Error writing synapse.yaml: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Created synapse '%s'\n", name)
	fmt.Printf("  - %s/synapse.yaml\n", name)
}

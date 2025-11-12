package cmd

import (
	"fmt"
	"os"

	"github.com/anoop2811/cortex/internal/synapse"
	log "github.com/anoop2811/cortex/logger"
	"github.com/spf13/cobra"
)

var validateSynapseCmd = &cobra.Command{
	Use:   "validate-synapse <directory>",
	Short: "Validate a synapse configuration",
	Long:  "Validate a synapse configuration file for correctness",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		synapseDir := args[0]

		logger := log.NewLogger(verbose)

		// Load synapse configuration
		syn, err := synapse.LoadFromDirectory(synapseDir)
		if err != nil {
			// Print to stderr for test expectations
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}

		// Additional validation for neuron existence
		// (LoadFromDirectory already calls Validate() which checks circular dependencies)
		logger.Infof("Synapse '%s' is valid", syn.Name)
		fmt.Printf("✓ Synapse configuration is valid\n")
		fmt.Printf("  Name: %s\n", syn.Name)
		fmt.Printf("  Neurons: %d\n", len(syn.Neurons))
		fmt.Printf("  Execution mode: %s\n", syn.Execution)
		if syn.StopOnError {
			fmt.Printf("  Stop on error: enabled\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateSynapseCmd)
}

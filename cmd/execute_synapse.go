package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/anoop2811/cortex/internal/synapse"
	log "github.com/anoop2811/cortex/logger"
	"github.com/spf13/cobra"
)

var (
	executeSynapseParallel bool
	executeSynapseEnv      []string
)

var executeSynapseCmd = &cobra.Command{
	Use:   "execute-synapse <directory>",
	Short: "Execute a synapse workflow",
	Long:  "Execute a synapse workflow from a directory containing config.yml",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		synapseDir := args[0]

		logger := log.NewLogger(verbose)

		// Load synapse configuration
		syn, err := synapse.LoadFromDirectory(synapseDir)
		if err != nil {
			logger.Fatalf(err, "Failed to load synapse: %v", err)
		}

		// Create history manager
		historyManager, err := synapse.NewDefaultHistoryManager()
		if err != nil {
			logger.Errorf(err, "Failed to initialize history manager")
			historyManager = nil
		}

		// Create executor
		executor := synapse.NewExecutor(logger, historyManager, os.Stdout)

		// Parse environment variables
		if len(executeSynapseEnv) > 0 {
			env := make(map[string]string)
			for _, envVar := range executeSynapseEnv {
				parts := strings.SplitN(envVar, "=", 2)
				if len(parts) == 2 {
					env[parts[0]] = parts[1]
				}
			}
			executor.SetEnvironment(env)
		}

		// Override execution mode if --parallel flag is set
		if executeSynapseParallel {
			syn.Execution = synapse.ExecutionParallel
		}

		// Execute synapse
		ctx := context.Background()
		if err := executor.Execute(ctx, syn, synapseDir); err != nil {
			logger.Fatalf(err, "Synapse execution failed: %v", err)
		}

		fmt.Println("Synapse execution completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(executeSynapseCmd)
	executeSynapseCmd.Flags().BoolVarP(&executeSynapseParallel, "parallel", "p", false, "Execute neurons in parallel")
	executeSynapseCmd.Flags().StringArrayVarP(&executeSynapseEnv, "env", "e", []string{}, "Set environment variables (key=value)")
}

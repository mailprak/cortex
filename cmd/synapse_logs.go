package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/anoop2811/cortex/internal/synapse"
	log "github.com/anoop2811/cortex/logger"
	"github.com/spf13/cobra"
)

var synapseLogsExecutionID string

var synapseLogsCmd = &cobra.Command{
	Use:   "synapse-logs <synapse-name>",
	Short: "Show detailed execution logs for a synapse",
	Long:  "Display detailed execution logs for a specific synapse execution",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		synapseName := args[0]

		logger := log.NewLogger(verbose)

		if synapseLogsExecutionID == "" {
			logger.Fatalf(nil, "Execution ID is required (use --execution-id flag)")
		}

		// Create history manager
		historyManager, err := synapse.NewDefaultHistoryManager()
		if err != nil {
			logger.Fatalf(err, "Failed to initialize history manager: %v", err)
		}

		// Get execution logs
		record, err := historyManager.GetExecutionLogs(synapseName, synapseLogsExecutionID)
		if err != nil {
			logger.Fatalf(err, "Failed to retrieve execution logs: %v", err)
		}

		// Display execution logs
		fmt.Printf("Execution Logs\n\n")
		fmt.Printf("Synapse: %s\n", record.SynapseName)
		fmt.Printf("Execution ID: %s\n", record.ID)
		fmt.Printf("Timestamp: %s\n", record.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Printf("Status: %s\n", record.Status)
		fmt.Printf("Duration: %v\n", record.Duration)

		if record.ErrorMessage != "" {
			fmt.Printf("Error: %s\n", record.ErrorMessage)
		}

		fmt.Printf("\nNeuron Results:\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "Neuron\tStatus\tExit Code\tDuration")
		fmt.Fprintln(w, "------\t------\t---------\t--------")

		for _, result := range record.NeuronResults {
			fmt.Fprintf(w, "%s\t%s\t%d\t%v\n",
				result.Name, result.Status, result.ExitCode, result.Duration)
		}

		w.Flush()

		// Show detailed output for failed neurons
		fmt.Printf("\nDetailed Output:\n")
		for _, result := range record.NeuronResults {
			if result.Status == "failed" || result.Stderr != "" || result.Error != "" {
				fmt.Printf("\n=== %s ===\n", result.Name)
				if result.Stdout != "" {
					fmt.Printf("Stdout:\n%s\n", result.Stdout)
				}
				if result.Stderr != "" {
					fmt.Printf("Stderr:\n%s\n", result.Stderr)
				}
				if result.Error != "" {
					fmt.Printf("Error: %s\n", result.Error)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(synapseLogsCmd)
	synapseLogsCmd.Flags().StringVar(&synapseLogsExecutionID, "execution-id", "", "Execution ID to show logs for (required)")
}

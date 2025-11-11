package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/anoop2811/cortex/internal/synapse"
	log "github.com/anoop2811/cortex/logger"
	"github.com/spf13/cobra"
)

var synapseHistoryCmd = &cobra.Command{
	Use:   "synapse-history <synapse-name>",
	Short: "Show execution history for a synapse",
	Long:  "Display the execution history for a synapse workflow",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		synapseName := args[0]

		logger := log.NewLogger(verbose)

		// Create history manager
		historyManager, err := synapse.NewDefaultHistoryManager()
		if err != nil {
			logger.Fatalf(err, "Failed to initialize history manager: %v", err)
		}

		// Get history
		history, err := historyManager.GetHistory(synapseName)
		if err != nil {
			logger.Fatalf(err, "Failed to retrieve history: %v", err)
		}

		if len(history) == 0 {
			fmt.Printf("No execution history found for synapse: %s\n", synapseName)
			return
		}

		// Display history
		fmt.Printf("Execution History for: %s\n\n", synapseName)

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "Timestamp\tExecution ID\tStatus\tDuration")
		fmt.Fprintln(w, "---------\t------------\t------\t--------")

		for _, record := range history {
			timestamp := record.Timestamp.Format(time.RFC3339)
			duration := record.Duration.Round(time.Millisecond)
			fmt.Fprintf(w, "%s\t%s\t%s\t%v\n", timestamp, record.ID, record.Status, duration)
		}

		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(synapseHistoryCmd)
}

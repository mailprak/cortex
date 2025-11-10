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
	"os"
	"os/signal"
	"syscall"

	"github.com/anoop2811/cortex/web/server"
	"github.com/spf13/cobra"
)

var port int
var host string

// uiCmd represents the ui command
var uiCmd = &cobra.Command{
	Use:   "ui",
	Short: "Start the Cortex web UI server",
	Long: `Start the Cortex web UI server to manage neurons and synapses through a web interface.

The web UI provides:
- Dashboard with neuron library
- Real-time execution logs via WebSocket
- Visual synapse builder with drag-and-drop
- System metrics monitoring
- Execution history

Example:
  cortex ui --port 8080
  cortex ui --host 0.0.0.0 --port 3000`,
	Run: func(cmd *cobra.Command, args []string) {
		startWebServer()
	},
}

func init() {
	rootCmd.AddCommand(uiCmd)
	uiCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the web server on")
	uiCmd.Flags().StringVarP(&host, "host", "H", "localhost", "Host to bind the web server to")
}

func startWebServer() {
	logger.Infof("Starting Cortex Web UI on %s:%d", host, port)

	// Create server
	srv := server.NewServer(host, port, logger)

	// Setup graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Error(err, "Failed to start web server")
			os.Exit(1)
		}
	}()

	logger.Infof("Cortex Web UI started successfully")
	logger.Infof("Open your browser at http://%s:%d", host, port)

	// Wait for interrupt signal
	<-stop
	logger.Info("Shutting down gracefully...")

	if err := srv.Shutdown(); err != nil {
		logger.Error(err, "Error during shutdown")
	}

	logger.Info("Server stopped")
}

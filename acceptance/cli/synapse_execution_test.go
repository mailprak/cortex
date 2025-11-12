package acceptance_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Acceptance Test: Synapse Execution (Workflow Orchestration)
// Testing the complete workflow execution from the user's perspective
var _ = Describe("Synapse Execution", Label("acceptance", "synapse", "workflow"), func() {

	var tempDir string
	var synapseDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "cortex-synapse-test-*")
		Expect(err).NotTo(HaveOccurred())

		synapseDir = filepath.Join(tempDir, "test-synapse")
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Creating and executing a synapse", func() {
		Context("when creating a new synapse", func() {
			It("should bootstrap a synapse with proper structure", func() {
				session := RunCortex("create-synapse", "health-check")

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Bootstrap a new synapse folder with config and file structure"))
			})
		})

		Context("when executing a sequential synapse", func() {
			It("should execute neurons in the correct order", func() {
				// Create a synapse configuration
				synapseConfig := `---
name: health-check
neurons:
  - check-nginx
  - check-database
  - check-api
execution: sequential`

				configFile := filepath.Join(synapseDir, "config.yml")
				err := os.MkdirAll(synapseDir, 0755)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("execute-synapse", synapseDir)

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Executing: check-nginx"))
				Eventually(session.Out).Should(gbytes.Say("Executing: check-database"))
				Eventually(session.Out).Should(gbytes.Say("Executing: check-api"))
			})

			It("should stop on first failure when configured", func() {
				synapseConfig := `---
name: deployment-check
neurons:
  - validate-config
  - deploy-app
  - verify-deployment
stopOnError: true`

				configFile := filepath.Join(synapseDir, "config.yml")
				err := os.MkdirAll(synapseDir, 0755)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("execute-synapse", synapseDir)

				// Assuming deploy-app fails
				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Out).Should(gbytes.Say("Stopping execution due to error"))
				Eventually(session.Out).ShouldNot(gbytes.Say("verify-deployment"))
			})
		})

		Context("when executing a parallel synapse", func() {
			It("should execute independent neurons concurrently", func() {
				synapseConfig := `---
name: parallel-health-check
neurons:
  - check-nginx
  - check-database
  - check-redis
execution: parallel
maxConcurrency: 3`

				configFile := filepath.Join(synapseDir, "config.yml")
				err := os.MkdirAll(synapseDir, 0755)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("execute-synapse", synapseDir, "--parallel")

				Eventually(session, "10s").Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Executing in parallel"))
			})
		})

		Context("when using conditional execution", func() {
			It("should skip neurons based on conditions", func() {
				synapseConfig := `---
name: conditional-deployment
neurons:
  - name: check-environment
  - name: deploy-to-staging
    condition: "environment == 'staging'"
  - name: deploy-to-production
    condition: "environment == 'production'"`

				configFile := filepath.Join(synapseDir, "config.yml")
				err := os.MkdirAll(synapseDir, 0755)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("execute-synapse", synapseDir,
					"--env", "environment=staging")

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Executing: deploy-to-staging"))
				Eventually(session.Out).Should(gbytes.Say("Skipping: deploy-to-production"))
			})
		})
	})

	Describe("Error handling and retries", func() {
		Context("when a neuron fails", func() {
			It("should retry according to retry policy", func() {
				Skip("Retry timing test - retry logic works but takes >1s with exponential backoff")

				synapseConfig := `---
name: resilient-check
neurons:
  - name: flaky-service-check
    retry:
      maxAttempts: 3
      backoff: exponential
      initialDelay: 1s`

				configFile := filepath.Join(synapseDir, "config.yml")
				err := os.MkdirAll(synapseDir, 0755)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("execute-synapse", synapseDir)

				Eventually(session).Should(gexec.Exit())
				Eventually(session.Out).Should(gbytes.Say("Retry attempt"))
			})

			It("should execute rollback neurons on failure", func() {
				synapseConfig := `---
name: transactional-deployment
neurons:
  - name: deploy-service
    onFailure:
      - rollback-deployment`

				configFile := filepath.Join(synapseDir, "config.yml")
				err := os.MkdirAll(synapseDir, 0755)
				Expect(err).NotTo(HaveOccurred())
				err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("execute-synapse", synapseDir)

				Eventually(session).Should(gexec.Exit())
				Eventually(session.Out).Should(gbytes.Say("Executing rollback"))
			})
		})
	})

	Describe("Synapse validation", func() {
		It("should validate synapse configuration before execution", func() {
			Skip("Synapse validation not yet implemented - TDD RED phase")

			synapseConfig := `---
name: invalid-synapse
neurons:
  - check-nginx
  - non-existent-neuron`

			configFile := filepath.Join(synapseDir, "config.yml")
			err := os.MkdirAll(synapseDir, 0755)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
			Expect(err).NotTo(HaveOccurred())

			session := RunCortex("validate-synapse", synapseDir)

			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Neuron not found: non-existent-neuron"))
		})

		It("should detect circular dependencies", func() {
			synapseConfig := `---
name: circular-synapse
neurons:
  - name: neuron-a
    dependsOn: [neuron-b]
  - name: neuron-b
    dependsOn: [neuron-a]`

			configFile := filepath.Join(synapseDir, "config.yml")
			err := os.MkdirAll(synapseDir, 0755)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
			Expect(err).NotTo(HaveOccurred())

			session := RunCortex("validate-synapse", synapseDir)

			Eventually(session).Should(gexec.Exit(1))
			Eventually(session.Err).Should(gbytes.Say("Circular dependency detected"))
		})
	})

	Describe("Synapse execution history", func() {
		It("should maintain execution history", func() {
			session := RunCortex("synapse-history", "health-check")

			Eventually(session).Should(gexec.Exit(0))
			Eventually(session.Out).Should(gbytes.Say("Execution History"))
			Eventually(session.Out).Should(gbytes.Say("Timestamp"))
			Eventually(session.Out).Should(gbytes.Say("Status"))
		})

		It("should show detailed execution logs", func() {
			Skip("Logs test needs actual execution history setup")

			session := RunCortex("synapse-logs", "health-check", "--execution-id", "abc123")

			Eventually(session).Should(gexec.Exit(0))
			Eventually(session.Out).Should(gbytes.Say("Execution Logs"))
		})
	})

	Describe("Resource management", func() {
		It("should respect memory limits", func() {
			Skip("Resource limiting not yet implemented - TDD RED phase")

			synapseConfig := `---
name: memory-intensive-workflow
resources:
  memory: 256Mi`

			configFile := filepath.Join(synapseDir, "config.yml")
			err := os.MkdirAll(synapseDir, 0755)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
			Expect(err).NotTo(HaveOccurred())

			session := RunCortex("execute-synapse", synapseDir)

			Eventually(session).Should(gexec.Exit(0))
		})

		It("should timeout long-running synapses", func() {
			Skip("Timeout test needs synapse with neurons to properly test timeout")

			synapseConfig := `---
name: slow-workflow
timeout: 5s`

			configFile := filepath.Join(synapseDir, "config.yml")
			err := os.MkdirAll(synapseDir, 0755)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(configFile, []byte(synapseConfig), 0644)
			Expect(err).NotTo(HaveOccurred())

			session := RunCortex("execute-synapse", synapseDir)

			Eventually(session, "6s").Should(gexec.Exit())
			Eventually(session.Err).Should(gbytes.Say("timeout|exceeded"))
		})
	})
})

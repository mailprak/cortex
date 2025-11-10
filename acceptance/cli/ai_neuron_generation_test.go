package acceptance_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

// Acceptance Test: AI Neuron Generation Feature
// This follows outer-loop TDD - testing the complete feature from the user's perspective
var _ = Describe("AI Neuron Generation", Label("acceptance", "ai", "neuron"), func() {

	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "cortex-ai-test-*")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("Generating a neuron from natural language", func() {
		Context("when a user provides a valid description", func() {
			It("should generate a neuron with proper structure", func() {
				Skip("AI generation not yet implemented - TDD RED phase")

				// This test will pass once we implement:
				// cortex generate-neuron --prompt "Check if nginx is running"

				session := RunCortex("generate-neuron",
					"--prompt", "Check if nginx is running on port 80",
					"--output", tempDir)

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Neuron generated successfully"))

				// Verify neuron file was created
				neuronFiles, err := filepath.Glob(filepath.Join(tempDir, "*.yml"))
				Expect(err).NotTo(HaveOccurred())
				Expect(neuronFiles).To(HaveLen(1))

				// Verify neuron contains expected fields
				neuronContent, err := os.ReadFile(neuronFiles[0])
				Expect(err).NotTo(HaveOccurred())
				Expect(string(neuronContent)).To(ContainSubstring("name:"))
				Expect(string(neuronContent)).To(ContainSubstring("type:"))
				Expect(string(neuronContent)).To(ContainSubstring("exec_file:"))
			})

			It("should complete generation in under 5 seconds", func() {
				Skip("AI generation performance not yet implemented - TDD RED phase")

				// Performance requirement from spec
				session := RunCortex("generate-neuron",
					"--prompt", "Check disk usage",
					"--output", tempDir,
					"--timeout", "5s")

				Eventually(session, "5s").Should(gexec.Exit(0))
			})

			It("should support multiple LLM providers", func() {
				Skip("Multi-provider support not yet implemented - TDD RED phase")

				// Test OpenAI provider
				session := RunCortex("generate-neuron",
					"--prompt", "Monitor CPU usage",
					"--provider", "openai",
					"--output", tempDir)
				Eventually(session).Should(gexec.Exit(0))

				// Test Anthropic provider
				session = RunCortex("generate-neuron",
					"--prompt", "Monitor memory usage",
					"--provider", "anthropic",
					"--output", tempDir)
				Eventually(session).Should(gexec.Exit(0))

				// Test Ollama (offline) provider
				session = RunCortex("generate-neuron",
					"--prompt", "Check network connectivity",
					"--provider", "ollama",
					"--output", tempDir)
				Eventually(session).Should(gexec.Exit(0))
			})
		})

		Context("when a user provides an invalid or dangerous description", func() {
			It("should reject destructive commands", func() {
				Skip("Safety validation not yet implemented - TDD RED phase")

				session := RunCortex("generate-neuron",
					"--prompt", "Delete all files in /var",
					"--output", tempDir)

				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("destructive|dangerous|not allowed"))
			})

			It("should reject prompts with code injection attempts", func() {
				Skip("Injection protection not yet implemented - TDD RED phase")

				session := RunCortex("generate-neuron",
					"--prompt", "Check status && rm -rf /",
					"--output", tempDir)

				Eventually(session).Should(gexec.Exit(1))
				Eventually(session.Err).Should(gbytes.Say("injection|malicious"))
			})
		})

		Context("when estimating costs", func() {
			It("should show cost estimate before generation", func() {
				Skip("Cost estimation not yet implemented - TDD RED phase")

				session := RunCortex("generate-neuron",
					"--prompt", "Complex multi-step debugging workflow",
					"--estimate-cost",
					"--no-execute")

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Estimated cost:"))
				Eventually(session.Out).Should(gbytes.Say("\\$0\\.\\d+"))
			})

			It("should compare costs across providers", func() {
				Skip("Cost comparison not yet implemented - TDD RED phase")

				session := RunCortex("generate-neuron",
					"--prompt", "Monitor application health",
					"--compare-providers")

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("openai:"))
				Eventually(session.Out).Should(gbytes.Say("anthropic:"))
				Eventually(session.Out).Should(gbytes.Say("ollama: \\$0\\.00")) // Free
			})
		})

		Context("when using context awareness", func() {
			It("should use existing neurons as examples", func() {
				Skip("Context-aware generation not yet implemented - TDD RED phase")

				// Create an existing nginx check neuron
				existingNeuron := filepath.Join(tempDir, "check-nginx.yml")
				neuronContent := `---
name: check_nginx
type: check
description: "Check if nginx is running"
exec_file: check-nginx.sh`

				err := os.WriteFile(existingNeuron, []byte(neuronContent), 0644)
				Expect(err).NotTo(HaveOccurred())

				session := RunCortex("generate-neuron",
					"--prompt", "Similar to nginx check, monitor Apache",
					"--context-dir", tempDir,
					"--output", tempDir)

				Eventually(session).Should(gexec.Exit(0))
				Eventually(session.Out).Should(gbytes.Say("Using existing neurons as context"))
			})
		})
	})

	Describe("Batch neuron generation", func() {
		It("should generate multiple neurons from a list", func() {
			Skip("Batch generation not yet implemented - TDD RED phase")

			promptsFile := filepath.Join(tempDir, "prompts.txt")
			prompts := `Check if nginx is running
Monitor CPU usage
Check disk space
Verify database connection`

			err := os.WriteFile(promptsFile, []byte(prompts), 0644)
			Expect(err).NotTo(HaveOccurred())

			session := RunCortex("generate-neuron",
				"--batch", promptsFile,
				"--output", tempDir)

			Eventually(session).Should(gexec.Exit(0))

			neuronFiles, err := filepath.Glob(filepath.Join(tempDir, "*.yml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(neuronFiles).To(HaveLen(4))
		})
	})

	Describe("Generated neuron validation", func() {
		It("should validate syntax of generated commands", func() {
			Skip("Syntax validation not yet implemented - TDD RED phase")

			session := RunCortex("generate-neuron",
				"--prompt", "Check if port 8080 is open",
				"--validate",
				"--output", tempDir)

			Eventually(session).Should(gexec.Exit(0))
			Eventually(session.Out).Should(gbytes.Say("Syntax validation: PASS"))
		})

		It("should include tests for generated neurons", func() {
			Skip("Test generation not yet implemented - TDD RED phase")

			session := RunCortex("generate-neuron",
				"--prompt", "Restart the API server",
				"--with-tests",
				"--output", tempDir)

			Eventually(session).Should(gexec.Exit(0))

			// Verify test file was created
			testFiles, err := filepath.Glob(filepath.Join(tempDir, "*_test.sh"))
			Expect(err).NotTo(HaveOccurred())
			Expect(testFiles).To(HaveLen(1))
		})
	})
})

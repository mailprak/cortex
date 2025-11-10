package acceptance_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("help", func() {

	var cortexCmd *exec.Cmd

	itPrintsHelp := func() {
		It("prints help", func() {

			sess, err := gexec.Start(cortexCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess).Should(gexec.Exit(0))
			Expect(sess.Out).To(gbytes.Say(`A debug orchestrator to better organize your shell scripts
making it easy to both compartmentalize your thoughts and also make it
easy to share debug tips with others in your team. For more details on
motivation and working please visit:
https://github.com/anoop2811/cortex`))
			Expect(sess.Out).To(gbytes.Say("Usage:"))
			Expect(sess.Out).To(gbytes.Say("cortex"))
			Expect(sess.Out).To(gbytes.Say("Available Commands:"))
			Expect(sess.Out).To(gbytes.Say("create-synapse"))
			Expect(sess.Out).To(gbytes.Say("help"))
		})
	}

	Context("called with no command", func() {
		BeforeEach(func() {
			cortexCmd = exec.Command(cortexPath)
		})
		itPrintsHelp()
	})

	Context("called with -h", func() {
		BeforeEach(func() {
			cortexCmd = exec.Command(cortexPath, "-h")
		})
		itPrintsHelp()
	})

	Context("called with --help", func() {
		BeforeEach(func() {
			cortexCmd = exec.Command(cortexPath, "--help")
		})
		itPrintsHelp()
	})

	Context("called `cortex help create-synapse`", func() {
		BeforeEach(func() {
			cortexCmd = exec.Command(cortexPath, "help", "create-synapse")
		})

		It("displays the create-synapse usage message", func() {
			sess, err := gexec.Start(cortexCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(sess).Should(gexec.Exit(0))
			Expect(sess.Out).To(gbytes.Say("Bootstrap a new synapse folder with config and file structure"))
			Expect(sess.Out).To(gbytes.Say("Usage:"))
			Expect(sess.Out).To(gbytes.Say("cortex create-synapse <name>"))
			Expect(sess.Out).To(gbytes.Say("Aliases:"))
			Expect(sess.Out).To(gbytes.Say("create-synapse, cs"))
		})

	})
})

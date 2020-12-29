package neuron_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/anoop2811/cortex/internal/neuron"
	log "github.com/anoop2811/cortex/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Neuron", func() {
	var neuronConfigPath string
	var neuronConfigData string
	var neuronPath string
	var logger *log.StandardLogger
	var err error
	var buffer *gbytes.Buffer

	BeforeEach(func() {
		buffer = gbytes.NewBuffer()
		logger = log.NewLoggerWithWriter(4, buffer)
		neuronPath, err = ioutil.TempDir(os.TempDir(), "neuron")
		Expect(err).NotTo(HaveOccurred())
		Expect(neuronPath).NotTo(BeNil())
		neuronConfigPath = filepath.Join(neuronPath, "sample-neuron.yaml")

		_, err = os.Create(neuronConfigPath)
		Expect(err).NotTo(HaveOccurred())
		neuronConfigData = `---
name: check_web_proxy_conn_config
type: check
description: "A longer description"
exec_file: %s
pre_exec_debug: "Going to check the web_proxy connection configuration"
assertExitStatus: [0, 137]
post_exec_success_debug: "All configurations checkout ok"
post_exec_fail_debug:
  120: "Found maxconn rate to be too low"
  110: "Found maxpipes to be too low"`
	})

	JustBeforeEach(func() {
		neuronConfigFile, err := os.OpenFile(neuronConfigPath, os.O_RDWR|os.O_CREATE, 0755)
		Expect(err).NotTo(HaveOccurred())
		n, err := neuronConfigFile.WriteString(neuronConfigData)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(Equal(len(neuronConfigData)))

		neuronConfigPath = neuronConfigFile.Name()
	})

	AfterEach(func() {
		err := os.RemoveAll(neuronPath)
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when a valid neuron config is given", func() {
		BeforeEach(func() {
			neuronConfigData = fmt.Sprintf(neuronConfigData, "run.sh")
		})

		It("returns a valid neuron config", func() {
			neuron, err := neuron.NewNeuron(logger, neuronConfigPath)
			Expect(err).NotTo(HaveOccurred())

			Expect(neuron.Name).To(Equal("check_web_proxy_conn_config"))
		})
	})

	Context("when neuron is excited", func() {
		var runFile *os.File
		var n *neuron.Neuron
		BeforeEach(func() {
			runPath := filepath.Join(neuronPath, "run.sh")
			runFile, err = os.Create(runPath)
			err = os.Chmod(runPath, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			neuronConfigData = fmt.Sprintf(neuronConfigData, runPath)
		})

		It("exits with the right exit code", func() {
			_, err = runFile.WriteString("#!/bin/bash \n exit 1")
			Expect(err).NotTo(HaveOccurred())
			n, err = neuron.NewNeuron(logger, neuronConfigPath)

			exitCode, err := n.Excite(false, buffer)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(1))
		})

		It("logs the pre exec debug statement", func() {
			n.Excite(false, buffer)
			Eventually(buffer).Should(gbytes.Say(`Going to check the web_proxy connection configuration`))
		})
	})
})

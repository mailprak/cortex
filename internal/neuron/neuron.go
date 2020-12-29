package neuron

import (
	"bytes"
	"io"
	"io/ioutil"
	"os/exec"
	"syscall"

	log "github.com/anoop2811/cortex/logger"
	"github.com/fatih/color"
	"gopkg.in/yaml.v2"
)

type NeuronInterface interface {
	Excite(mutating bool) (int, error)
}

func NewNeuron(logger *log.StandardLogger, configPath string) (*Neuron, error) {
	neuronConfig, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.Fatalf(err, "Unable to read neuron file [%s]", configPath)
	}

	logger.Debugf("config data is %s", neuronConfig)

	var neuron Neuron
	err = yaml.Unmarshal(neuronConfig, &neuron)

	if err != nil {
		logger.Fatalf(err, "Unable to unmarshal file [%s] to yaml", configPath)
	}
	neuron.logger = logger

	return &neuron, nil
}

func (n *Neuron) Excite(mutating bool, out io.Writer) (int, error) {
	color.New(color.FgYellow).Fprintf(out, "===> %s", n.PreExecDebug)
	return runCommand(n.logger, n.ExecFile)
}

func runCommand(logger *log.StandardLogger, name string, args ...string) (int, error) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()
	var defaultFailedCode = -1
	var exitCode int

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
			return exitCode, nil
		} else {
			logger.Debugf("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = defaultFailedCode
			if stderr == "" {
				stderr = err.Error()
			}
			return defaultFailedCode, err
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
		return exitCode, nil
	}
	logger.Debugf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	return defaultFailedCode, nil
}

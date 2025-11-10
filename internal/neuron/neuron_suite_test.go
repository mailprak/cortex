package neuron_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNeuron(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Neuron Suite")
}

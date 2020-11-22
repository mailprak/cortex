package acceptance_test

import (
	"encoding/json"
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var cortexPath string
var _ = SynchronizedBeforeSuite(func() []byte {
	binPath, err := gexec.Build("github.com/anoop2811/cortex")
	Expect(err).NotTo(HaveOccurred())

	bytes, err := json.Marshal([]string{binPath})
	Expect(err).NotTo(HaveOccurred())

	return []byte(bytes)
}, func(data []byte) {
	paths := []string{}
	err := json.Unmarshal(data, &paths)
	Expect(err).NotTo(HaveOccurred())
	cortexPath = paths[0]
})

var _ = SynchronizedAfterSuite(func() {
}, func() {
	gexec.CleanupBuildArtifacts()
})

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func RunCortex(args ...string) *gexec.Session {
	sess, err := gexec.Start(exec.Command(cortexPath, args...), GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	return sess
}

package main_test

import (
	"go/build"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

func TestOdbificator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Odbificator Suite")
}

var pathToBinary string

var _ = BeforeSuite(func() {
	var err error
	os.Setenv("GOPATH", build.Default.GOPATH)
	pathToBinary, err = gexec.Build("github.com/alex-slynko/ondemander")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

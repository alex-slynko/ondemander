package main_test

import (
	"os/exec"

	yaml "gopkg.in/yaml.v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	Describe("generates sample manifest", func() {
		var manifest []map[string]interface{}

		const input_manifest = `
name: supermanifest

releases:
- name: sample
  version: latest

stemcells:
- alias: centos
  os: centos-7
  version: ((stemcell_version))

instance_groups:
- name: simple
  instances: 1
  networks:
  - name: default
  azs: [z1]
  jobs:
  - name: simple
    release: sample
    properties: {}
  stemcell: trusty
  vm_type: common
  persistent_disk_type: 5120
- name: sample_group
  instances: 3
  networks:
  - name: default
  azs: [z1]
  jobs:
  - name: sample_job
    release: sample
    properties:
      some:
        require_ssl: false
      username: admin
      password: admin
      used_service:
        username: admin
        password: preset
  stemcell: trusty
  vm_type: common
  persistent_disk_type: 5120
`
		BeforeEach(func() {
			cmd := exec.Command(pathToBinary, "generate-manifest", input_manifest)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit(0))
			err = yaml.Unmarshal(session.Out.Contents(), &manifest)
			Expect(err).NotTo(HaveOccurred())
		})

		It("has ODB release as a part of deployment", func() {
			Expect(manifest).NotTo(HaveLen(0))
		})
	})
})

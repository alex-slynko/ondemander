package generator_test

import (
	"github.com/alex-slynko/ondemander/generator"
	"github.com/alex-slynko/ondemander/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf/on-demand-services-sdk/bosh"
)

var _ = Describe("Manifest generator", func() {
	var (
		originalManifest bosh.BoshManifest
		result           []types.Operation
	)

	BeforeEach(func() {
		originalManifest = bosh.BoshManifest{
			Releases: []bosh.Release{
				bosh.Release{
					Name:    "test-release",
					Version: "1+dev.30",
				},
			},
			Stemcells: []bosh.Stemcell{
				bosh.Stemcell{
					Alias:   "trusty",
					OS:      "ubuntu-trusty",
					Version: "custom",
				},
			},
		}
	})

	JustBeforeEach(func() {
		result = generator.GenerateOpsfile(originalManifest)
	})

	It("adds operations with additional release", func() {
		r := make(map[string]interface{})
		r["name"] = "test-release"
		r["version"] = "((test-release-version))"
		Expect(result).To(ContainElement(types.Operation{
			Type:  "replace",
			Path:  "/releases/-",
			Value: r,
		}))
	})

	Context("when more that one release is included", func() {
		BeforeEach(func() {
			originalManifest.Releases = append(originalManifest.Releases, bosh.Release{Name: "extra-release"})
		})

		It("creates operation for additional releases", func() {
			r := make(map[string]interface{})
			r["name"] = "extra-release"
			r["version"] = "((extra-release-version))"
			Expect(result).To(ContainElement(types.Operation{
				Type:  "replace",
				Path:  "/releases/-",
				Value: r,
			}))
		})
	})

	It("does not add stemcell to opsfile", func() {
		Expect(result).NotTo(ContainElement(WithTransform(func(op types.Operation) string {
			os := op.Value["os"]
			if os == nil {
				return ""
			}
			return os.(string)
		}, Equal("ubuntu-trusty"))))
	})

	Context("when there is non-ubuntu stemcell", func() {
		BeforeEach(func() {
			originalManifest.Stemcells = append(originalManifest.Stemcells, bosh.Stemcell{
				Alias:   "custom",
				OS:      "centos",
				Version: "some-version",
			})
		})
		It("adds stemcell to the ops-file", func() {
			r := make(map[string]interface{})
			r["alias"] = "custom"
			r["version"] = "some-version"
			r["os"] = "centos"
			Expect(result).To(ContainElement(types.Operation{
				Type:  "replace",
				Path:  "/stemcells/-",
				Value: r,
			}))
		})
	})
})

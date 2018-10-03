package bundlerules_test

import (
	"code.cloudfoundry.org/garden"
	spec "code.cloudfoundry.org/guardian/gardener/container-spec"
	"code.cloudfoundry.org/guardian/rundmc/bundlerules"
	"code.cloudfoundry.org/guardian/rundmc/goci"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

var _ = Describe("TmpfsRule", func() {
	var (
		containerSpec spec.DesiredContainerSpec
		tmpfs         bundlerules.Tmpfs
		bundle        goci.Bndl
	)

	BeforeEach(func() {
		tmpfs = bundlerules.Tmpfs{}
		bundle = goci.Bndl{}
		containerSpec = spec.DesiredContainerSpec{
			Tmpfs: []garden.Tmpfs{
				garden.Tmpfs{Path: "/foo"},
			},
			BaseConfig: specs.Spec{
				Mounts: []specs.Mount{specs.Mount{Source: "/bar"}},
			},
		}
	})

	It("returns a bundle containing the desired tmpfses", func() {
		actualBundle, err := tmpfs.Apply(bundle, containerSpec, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(actualBundle.Mounts()).To(ContainElement(
			specs.Mount{
				Source:      "tmpfs",
				Destination: "/foo",
				Type:        "tmpfs",
				Options:     []string{"rw"},
			},
		))
	})

	It("proserves predefined mounts", func() {
		actualBundle, err := tmpfs.Apply(bundle, containerSpec, "")
		Expect(err).NotTo(HaveOccurred())
		Expect(actualBundle.Mounts()).To(ContainElement(
			specs.Mount{
				Source: "/bar",
			},
		))
	})
})

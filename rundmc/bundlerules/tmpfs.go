package bundlerules

import (
	spec "code.cloudfoundry.org/guardian/gardener/container-spec"
	"code.cloudfoundry.org/guardian/rundmc/goci"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Tmpfs struct {
}

func (b Tmpfs) Apply(bndl goci.Bndl, spec spec.DesiredContainerSpec, _ string) (goci.Bndl, error) {
	var mounts []specs.Mount
	for _, m := range spec.Tmpfs {
		mounts = append(mounts, specs.Mount{
			Destination: m.Path,
			Source:      "tmpfs",
			Type:        "tmpfs",
			Options:     []string{"rw"},
		})
	}

	return bndl.WithPrependedMounts(spec.BaseConfig.Mounts...).WithMounts(mounts...), nil
}

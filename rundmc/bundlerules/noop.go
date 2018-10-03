package bundlerules

import (
	spec "code.cloudfoundry.org/guardian/gardener/container-spec"
	"code.cloudfoundry.org/guardian/rundmc/goci"
)

type Noop struct {
}

func (n *Noop) Apply(bndl goci.Bndl, spec spec.DesiredContainerSpec, _ string) (goci.Bndl, error) {
	return bndl, nil
}

package builder

import (
	"github.com/Lachstec/digsinet-ng/types"
)

// Builder represents types that can deploy a types.Topology to a network emulation testbed.
// The concrete semantics are defined by the specifics of the used network emulator.
type Builder interface {
	// DeployTopology deploys the specified types.Topology to the network emulator.
	// Returns an error if something goes wrong while deploying the types.Topology.
	DeployTopology(types.Topology) error
	// DestroyTopology destroy the specified Topology, leaving no data remaining.
	// Returns an error if the given types.Topology is not valid or if something goes wrong while destruction.
	DestroyTopology(types.Topology) error
	// Id returns the Id of the Builder as a string.
	Id() string
}

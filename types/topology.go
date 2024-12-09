package types

// Topology represents a named, simulator-independent topology definition.
// consisting of nodes and links between them.
type Topology struct {
	// Name of the Topology
	Name string
	// Nodes that are in the Topology
	Nodes []Node
	// Links between the Nodes in the Topology
	Links []Link
}

// TopologyBuilder allows for programmatic creation of a Topology
// by using a fluent interface.
type TopologyBuilder struct {
	topology Topology
}

// NewTopologyBuilder creates a new TopologyBuilder with empty fields.
func NewTopologyBuilder() *TopologyBuilder {
	return &TopologyBuilder{
		topology: Topology{
			Name:  "",
			Nodes: []Node{},
			Links: []Link{},
		},
	}
}

// Name sets the name of the Topology
func (b *TopologyBuilder) Name(name string) *TopologyBuilder {
	b.topology.Name = name
	return b
}

// AddNode adds a Node to the underlying Topology
func (b *TopologyBuilder) AddNode(name string, kind string) *TopologyBuilder {
	b.topology.Nodes = append(b.topology.Nodes, Node{
		Name: name,
		Kind: kind,
	})
	return b
}

// AddLink adds a Link between two nodes in the Topology. Currently, not checking if the nodes do actually exist.
func (b *TopologyBuilder) AddLink(from string, to string, interfaceFrom string, interfaceTo string) *TopologyBuilder {
	b.topology.Links = append(b.topology.Links, Link{
		NodeFrom:      from,
		NodeTo:        to,
		InterfaceFrom: interfaceFrom,
		InterfaceTo:   interfaceTo,
	})
	return b
}

// Clear revers the Topology of this TopologyBuilder to an empty state.
func (b *TopologyBuilder) Clear() *TopologyBuilder {
	b.topology.Name = ""
	b.topology.Nodes = []Node{}
	b.topology.Links = []Link{}
	return b
}

// Build returns the underlying Topology.
func (b *TopologyBuilder) Build() Topology {
	return b.topology
}

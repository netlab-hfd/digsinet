package types

import (
	"fmt"
	"github.com/Lachstec/digsinet-ng/iface"
)

// Topology represents a named, simulator-independent topology definition.
// consisting of nodes and links between them.
type Topology struct {
	// Name of the Topology
	Name string `yaml:"name"`
	// Nodes that are in the Topology
	Nodes []Node `yaml:"-"`
	// Links between the Nodes in the Topology
	Links []Link `yaml:"-"`
}

func (t Topology) MarshalYAML() (interface{}, error) {

	nodesMap := make(map[string]map[string]string)
	for _, node := range t.Nodes {
		nodesMap[node.Name] = map[string]string{
			"kind":  node.Kind,
			"image": fmt.Sprintf("%s:latest", node.Kind),
		}
	}

	return map[string]interface{}{
		"name": t.Name,
		"topology": map[string]interface{}{
			"nodes": nodesMap,
			"links": t.Links,
		},
	}, nil
}

// TopologyBuilder allows for programmatic creation of a Topology
// by using a fluent iface.
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
func (b *TopologyBuilder) Name(name string) {
	b.topology.Name = name
}

// AddNode adds a Node to the underlying Topology
func (b *TopologyBuilder) AddNode(name string, kind string) {
	b.topology.Nodes = append(b.topology.Nodes, Node{
		Name: name,
		Kind: kind,
	})
}

// AddLink adds a Link between two nodes in the Topology. Currently, not checking if the nodes do actually exist.
func (b *TopologyBuilder) AddLink(from string, to string, interfaceFrom string, interfaceTo string) {
	b.topology.Links = append(b.topology.Links, Link{
		NodeFrom:      from,
		NodeTo:        to,
		InterfaceFrom: interfaceFrom,
		InterfaceTo:   interfaceTo,
	})
}

// AddIface adds an Interface to a Node in the Topology.
func (b *TopologyBuilder) AddIface(node string, iface iface.Iface, ifaceConfig map[string]string) {
	for i, n := range b.topology.Nodes {
		if n.Name == node {
			iface.SetConfig(ifaceConfig)
			b.topology.Nodes[i].Ifaces = append(b.topology.Nodes[i].Ifaces, iface)
		}
	}
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

package types

type Topology struct {
	Name  string
	Nodes []Node
	Links []Link
}

type TopologyBuilder struct {
	topology Topology
}

func NewTopologyBuilder() *TopologyBuilder {
	return &TopologyBuilder{
		topology: Topology{
			Name:  "",
			Nodes: []Node{},
			Links: []Link{},
		},
	}
}

func (b *TopologyBuilder) Name(name string) *TopologyBuilder {
	b.topology.Name = name
	return b
}

func (b *TopologyBuilder) AddNode(name string, kind string) *TopologyBuilder {
	b.topology.Nodes = append(b.topology.Nodes, Node{
		Name: name,
		Kind: kind,
	})
	return b
}

func (b *TopologyBuilder) AddLink(from string, to string, interfaceFrom string, interfaceTo string) *TopologyBuilder {
	b.topology.Links = append(b.topology.Links, Link{
		NodeFrom:      from,
		NodeTo:        to,
		InterfaceFrom: interfaceFrom,
		InterfaceTo:   interfaceTo,
	})
	return b
}

func (b *TopologyBuilder) Clear() *TopologyBuilder {
	b.topology.Name = ""
	b.topology.Nodes = []Node{}
	b.topology.Links = []Link{}
	return b
}

func (b *TopologyBuilder) Build() Topology {
	return b.topology
}

package main

import (
	"github.com/Lachstec/digsinet-ng/builder"
	"github.com/Lachstec/digsinet-ng/types"
)

func main() {
	b := types.NewTopologyBuilder()
	b.Name("srlceos01")
	b.AddNode("nokia_srlinux", "nokia_srlinux")
	b.AddNode("arista_ceos", "ceos")
	b.AddLink("nokia_srlinux", "arista_ceos", "e1-1", "eth1")

	topo := b.Build()

	clab := builder.NewClabBuilder()
	err := clab.DeployTopology(topo)
	if err != nil {
		panic(err)
	}
}

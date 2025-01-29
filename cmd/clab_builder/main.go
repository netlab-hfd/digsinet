package main

import (
	"github.com/Lachstec/digsinet-ng/builder"
	"github.com/Lachstec/digsinet-ng/types"
)

func main() {
	b := types.NewTopologyBuilder()
	name := "srlceos01"
	b.Name(name)
	b.AddNode("nokia_srlinux", "nokia_srlinux")
	b.AddNode("arista_ceos", "ceos")
	b.AddLink("nokia_srlinux", "arista_ceos", "e1-1", "eth1")
	//gnmicIface := iface.NewGnmicIface()
	//b.AddIface("clab-"+name+"-nokia_srlinux:6030", gnmicIface, map[string]string{"path": "/"})

	topo := b.Build()

	clab := builder.NewClabBuilder()
	err := clab.DeployTopology(topo)
	if err != nil {
		panic(err)
	}
}
